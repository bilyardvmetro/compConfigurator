-- Функция: проверка совместимости сборки по её id
CREATE OR REPLACE FUNCTION check_assembly_compatibility(p_assembly_id BIGINT)
    RETURNS VOID AS
$$
DECLARE
    v_cpu                 cpus%ROWTYPE;
    v_mb                  motherboards%ROWTYPE;
    v_psu                 psus%ROWTYPE;
    v_case                cases%ROWTYPE;
    v_cooler              cpu_coolers%ROWTYPE;
    v_gpu                 gpus%ROWTYPE;
    v_total_ram_modules   INTEGER := 0;
    v_total_ram_capacity  INTEGER := 0;
    v_max_ram_freq        INTEGER := 0;
    v_nvme_count          INTEGER := 0;
    v_sata_count          INTEGER := 0;
    v_bays_35_used        INTEGER := 0;
    v_bays_25_used        INTEGER := 0;
    v_required_power_watt INTEGER := 0;
BEGIN
    ----------------------------------------------------------------
    -- 1. Загрузить основную информацию о сборке и компонентах
    ----------------------------------------------------------------
    SELECT c.*, m.*, p.*, ca.*, co.*, g.*
    INTO v_cpu, v_mb, v_psu, v_case, v_cooler, v_gpu
    FROM assemblies a
             JOIN cpus c ON c.cpu_id = a.cpu_id
             JOIN motherboards m ON m.motherboard_id = a.motherboard_id
             JOIN psus p ON p.psu_id = a.psu_id
             JOIN cases ca ON ca.case_id = a.case_id
             LEFT JOIN cpu_coolers co ON co.cooler_id = a.cooler_id
             LEFT JOIN gpus g ON g.gpu_id = a.gpu_id
    WHERE a.assembly_id = p_assembly_id;

    IF NOT FOUND THEN
        -- Если сборки нет (например, удалена) — тихо выходим
        RETURN;
    END IF;

    ----------------------------------------------------------------
    -- 2. CPU ↔ Motherboard: socket
    ----------------------------------------------------------------
    IF v_cpu.socket_code <> v_mb.socket_code THEN
        RAISE EXCEPTION
            'Несовместимость сокетов: CPU % (%) и Motherboard % (%)',
            v_cpu.name, v_cpu.socket_code,
            v_mb.name, v_mb.socket_code
            USING ERRCODE = 'check_violation';
    END IF;

    ----------------------------------------------------------------
    -- 3. Motherboard ↔ Case: форм-фактор
    ----------------------------------------------------------------
    IF NOT EXISTS (SELECT 1
                   FROM case_form_factor_support cffs
                   WHERE cffs.case_id = v_case.case_id
                     AND cffs.form_factor_code = v_mb.form_factor_code) THEN
        RAISE EXCEPTION
            'Форм-фактор материнской платы % не поддерживается корпусом %',
            v_mb.form_factor_code, v_case.name
            USING ERRCODE = 'check_violation';
    END IF;

    ----------------------------------------------------------------
    -- 4. Cooler ↔ CPU и Cooler ↔ Case
    ----------------------------------------------------------------
    IF v_cooler.cooler_id IS NOT NULL THEN
        -- 4.1. Кулер должен поддерживать сокет CPU
        IF NOT EXISTS (SELECT 1
                       FROM cooler_sockets cs
                       WHERE cs.cooler_id = v_cooler.cooler_id
                         AND cs.socket_code = v_cpu.socket_code) THEN
            RAISE EXCEPTION
                'Кулер % не поддерживает сокет CPU %',
                v_cooler.name, v_cpu.socket_code
                USING ERRCODE = 'check_violation';
        END IF;

        -- 4.2. TDP кулера должен покрывать TDP процессора
        IF v_cooler.tdp_limit_watt IS NOT NULL
            AND v_cpu.tdp_watt IS NOT NULL
            AND v_cooler.tdp_limit_watt < v_cpu.tdp_watt THEN
            RAISE EXCEPTION
                'Кулер % не выдерживает TDP процессора (% Вт < % Вт)',
                v_cooler.name, v_cooler.tdp_limit_watt, v_cpu.tdp_watt
                USING ERRCODE = 'check_violation';
        END IF;

        -- 4.3. Высота кулера должна помещаться в корпус
        IF v_cooler.height_mm IS NOT NULL
            AND v_case.max_cooler_height_mm IS NOT NULL
            AND v_cooler.height_mm > v_case.max_cooler_height_mm THEN
            RAISE EXCEPTION
                'Кулер % слишком высокий (% мм) для корпуса % (макс % мм)',
                v_cooler.name, v_cooler.height_mm,
                v_case.name, v_case.max_cooler_height_mm
                USING ERRCODE = 'check_violation';
        END IF;
    END IF;

    ----------------------------------------------------------------
    -- 5. GPU ↔ Case (габариты)
    ----------------------------------------------------------------
    IF v_gpu.gpu_id IS NOT NULL THEN
        -- длина
        IF v_gpu.length_mm IS NOT NULL
            AND v_case.max_gpu_length_mm IS NOT NULL
            AND v_gpu.length_mm > v_case.max_gpu_length_mm THEN
            RAISE EXCEPTION
                'Видеокарта % слишком длинная (% мм) для корпуса % (макс % мм)',
                v_gpu.name, v_gpu.length_mm,
                v_case.name, v_case.max_gpu_length_mm
                USING ERRCODE = 'check_violation';
        END IF;

        -- толщина в слотах
        IF v_gpu.width_slots IS NOT NULL
            AND v_case.max_gpu_width_slots IS NOT NULL
            AND v_gpu.width_slots > v_case.max_gpu_width_slots THEN
            RAISE EXCEPTION
                'Видеокарта % слишком толстая (% слота) для корпуса % (макс %)',
                v_gpu.name, v_gpu.width_slots,
                v_case.name, v_case.max_gpu_width_slots
                USING ERRCODE = 'check_violation';
        END IF;
    END IF;

    ----------------------------------------------------------------
    -- 6. PSU ↔ Case (форм-фактор БП)
    ----------------------------------------------------------------
    IF v_psu.form_factor IS NOT NULL
        AND v_case.psu_form_factor IS NOT NULL
        AND v_psu.form_factor <> v_case.psu_form_factor THEN
        RAISE EXCEPTION
            'Форм-фактор БП % не подходит к корпусу % (ожидалось %)',
            v_psu.form_factor, v_case.name, v_case.psu_form_factor
            USING ERRCODE = 'check_violation';
    END IF;

    -- длина БП
    IF v_psu.length_mm IS NOT NULL
        AND v_case.max_psu_length_mm IS NOT NULL
        AND v_psu.length_mm > v_case.max_psu_length_mm THEN
        RAISE EXCEPTION
            'БП % слишком длинный (% мм) для корпуса % (макс % мм)',
            v_psu.name, v_psu.length_mm,
            v_case.name, v_case.max_psu_length_mm
            USING ERRCODE = 'check_violation';
    END IF;

    ----------------------------------------------------------------
    -- 7. RAM ↔ Motherboard ↔ CPU
    ----------------------------------------------------------------
    SELECT COALESCE(SUM(k.module_count), 0),
           COALESCE(SUM(k.total_capacity_gb), 0),
           COALESCE(MAX(k.freq_mhz), 0)
    INTO
        v_total_ram_modules,
        v_total_ram_capacity,
        v_max_ram_freq
    FROM assembly_ram_kits ark
             JOIN ram_kits k ON k.ram_kit_id = ark.ram_kit_id
    WHERE ark.assembly_id = p_assembly_id;

    -- Тип RAM должен совпадать с материнкой
    IF EXISTS (SELECT 1
               FROM assembly_ram_kits ark
                        JOIN ram_kits k ON k.ram_kit_id = ark.ram_kit_id
               WHERE ark.assembly_id = p_assembly_id
                 AND (k.ram_type IS NOT NULL AND v_mb.ram_type IS NOT NULL)
                 AND k.ram_type <> v_mb.ram_type) THEN
        RAISE EXCEPTION
            'Тип RAM в сборке не совпадает с типом RAM материнской платы %',
            v_mb.ram_type
            USING ERRCODE = 'check_violation';
    END IF;

    -- Максимальная частота RAM должна быть не выше, чем поддерживают CPU и MB
    IF v_max_ram_freq > 0 THEN
        IF v_cpu.supported_ram_freq_max_mhz IS NOT NULL
            AND v_mb.ram_freq_max_mhz IS NOT NULL
            AND v_max_ram_freq >
                LEAST(v_cpu.supported_ram_freq_max_mhz, v_mb.ram_freq_max_mhz) THEN
            RAISE EXCEPTION
                'Частота RAM (%) выше, чем поддерживают CPU (% МГц) и MB (% МГц)',
                v_max_ram_freq,
                v_cpu.supported_ram_freq_max_mhz,
                v_mb.ram_freq_max_mhz
                USING ERRCODE = 'check_violation';
        END IF;
    END IF;

    -- Количество модулей не должно превышать количество слотов
    IF v_total_ram_modules > v_mb.ram_slots THEN
        RAISE EXCEPTION
            'Слишком много модулей RAM (%), материнская плата % поддерживает максимум %',
            v_total_ram_modules, v_mb.name, v_mb.ram_slots
            USING ERRCODE = 'check_violation';
    END IF;

    -- Общий объём не должен превышать максимум MB
    IF v_total_ram_capacity > v_mb.ram_capacity_max_gb THEN
        RAISE EXCEPTION
            'Слишком большой общий объём RAM (% ГБ), материнская плата % поддерживает максимум % ГБ',
            v_total_ram_capacity, v_mb.name, v_mb.ram_capacity_max_gb
            USING ERRCODE = 'check_violation';
    END IF;

    ----------------------------------------------------------------
    -- 8. Drives ↔ Motherboard ↔ Case
    ----------------------------------------------------------------
    SELECT COALESCE(SUM(CASE WHEN d.interface = 'NVMe' THEN 1 ELSE 0 END), 0)  AS nvme_count,
           COALESCE(SUM(CASE WHEN d.interface = 'SATA' THEN 1 ELSE 0 END), 0)  AS sata_count,
           COALESCE(SUM(CASE WHEN d.form_factor = '3.5' THEN 1 ELSE 0 END), 0) AS bays_35_used,
           COALESCE(SUM(CASE WHEN d.form_factor = '2.5' THEN 1 ELSE 0 END), 0) AS bays_25_used
    INTO
        v_nvme_count, v_sata_count, v_bays_35_used, v_bays_25_used
    FROM assembly_drives ad
             JOIN storage_drives d ON d.drive_id = ad.drive_id
    WHERE ad.assembly_id = p_assembly_id;

    -- NVMe <= M.2 слоты
    IF v_nvme_count > v_mb.m2_slots_count THEN
        RAISE EXCEPTION
            'Слишком много NVMe-накопителей (%), материнская плата % имеет только % слотов M.2',
            v_nvme_count, v_mb.name, v_mb.m2_slots_count
            USING ERRCODE = 'check_violation';
    END IF;

    -- SATA <= SATA порты
    IF v_sata_count > v_mb.sata_ports_count THEN
        RAISE EXCEPTION
            'Слишком много SATA-накопителей (%), материнская плата % имеет только % SATA портов',
            v_sata_count, v_mb.name, v_mb.sata_ports_count
            USING ERRCODE = 'check_violation';
    END IF;

    -- Корзины в корпусе
    IF v_bays_35_used > v_case.drive_bays_3_5_count THEN
        RAISE EXCEPTION
            'Слишком много 3.5" дисков (%), корпус % поддерживает только % корзин 3.5"',
            v_bays_35_used, v_case.name, v_case.drive_bays_3_5_count
            USING ERRCODE = 'check_violation';
    END IF;

    IF v_bays_25_used > v_case.drive_bays_2_5_count THEN
        RAISE EXCEPTION
            'Слишком много 2.5" дисков (%), корпус % поддерживает только % корзин 2.5"',
            v_bays_25_used, v_case.name, v_case.drive_bays_2_5_count
            USING ERRCODE = 'check_violation';
    END IF;

    ----------------------------------------------------------------
    -- 9. PSU ↔ суммарный TDP
    -- (упрощённая модель: CPU + GPU + небольшой запас)
    ----------------------------------------------------------------
    v_required_power_watt :=
            COALESCE(v_cpu.tdp_watt, 0)
                + COALESCE(v_gpu.tdp_watt, 0)
                + 100; -- запас на остальное железо

    IF v_psu.power_watt IS NOT NULL
        AND v_psu.power_watt < v_required_power_watt THEN
        RAISE EXCEPTION
            'Мощности БП % (% Вт) недостаточно, требуется примерно ≥ % Вт',
            v_psu.name, v_psu.power_watt, v_required_power_watt
            USING ERRCODE = 'check_violation';
    END IF;

    RETURN;
END;
$$ LANGUAGE plpgsql;