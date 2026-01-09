SET search_path TO pc_configurator;

-- Функция: проверка совместимости сборки по её id
CREATE OR REPLACE FUNCTION pc_configurator.check_assembly_compatibility(p_assembly_id BIGINT)
    RETURNS VOID
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
DECLARE
    v_asm                 assemblies%ROWTYPE;

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
    SELECT *
    INTO v_asm
    FROM assemblies
    WHERE assembly_id = p_assembly_id;

    IF NOT FOUND THEN
        -- Если сборки нет (например, удалена) — тихо выходим
        RETURN;
    END IF;

    -- Обязательные компоненты
    SELECT * INTO v_cpu  FROM cpus         WHERE cpu_id = v_asm.cpu_id;
    SELECT * INTO v_mb   FROM motherboards WHERE motherboard_id = v_asm.motherboard_id;
    SELECT * INTO v_psu  FROM psus         WHERE psu_id = v_asm.psu_id;
    SELECT * INTO v_case FROM cases        WHERE case_id = v_asm.case_id;

    -- Опциональные компоненты
    IF v_asm.cooler_id IS NOT NULL THEN
        SELECT * INTO v_cooler FROM cpu_coolers WHERE cooler_id = v_asm.cooler_id;
    ELSE
        v_cooler.cooler_id := NULL;
    END IF;

    IF v_asm.gpu_id IS NOT NULL THEN
        SELECT * INTO v_gpu FROM gpus WHERE gpu_id = v_asm.gpu_id;
    ELSE
        v_gpu.gpu_id := NULL;
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
                'Видеокарта % слишком толстая (%.1f слота) для корпуса % (макс %.1f)',
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
$$;

------------------------------------------------------------
-- 2) МЯГКАЯ ПРОВЕРКА ДЛЯ UI/API (не кидает, а возвращает ok + message)
------------------------------------------------------------

-- Функция: вернуть признак совместимости и текст ошибки (если есть)
CREATE OR REPLACE FUNCTION pc_configurator.is_assembly_compatible(
    p_assembly_id BIGINT,
    OUT is_ok BOOLEAN,
    OUT error_message TEXT
)
    RETURNS RECORD
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
BEGIN
    -- Пытаемся вызвать "строгую" проверку, которая кидает исключения
    PERFORM pc_configurator.check_assembly_compatibility(p_assembly_id);

    -- Если исключения не было — всё ок
    is_ok := TRUE;
    error_message := NULL;
EXCEPTION
    WHEN check_violation THEN -- наши бизнес-ошибки совместимости
        is_ok := FALSE;
        error_message := SQLERRM;
    WHEN OTHERS THEN -- системные ошибки (например, FK)
        RAISE; -- их пробрасываем дальше
END;
$$;

-- SELECT * FROM is_assembly_compatible(1);
-- -- вернёт, например:  is_ok = true, error_message = null

------------------------------------------------------------
-- 3) ЦЕНЫ
------------------------------------------------------------

-- Вспомогательная функция: минимальная цена для конкретного компонента
CREATE OR REPLACE FUNCTION pc_configurator.get_min_price(
    p_component_type VARCHAR,
    p_component_id BIGINT
)
    RETURNS NUMERIC
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
DECLARE
    v_price NUMERIC;
BEGIN
    SELECT MIN(price)
    INTO v_price
    FROM pc_configurator.product_offers
    WHERE component_type = p_component_type
      AND component_id = p_component_id
      AND available = TRUE;

    RETURN v_price; -- может быть NULL, если предложений нет
END;
$$;

-- Пересчитать общую цену сборки и записать в assemblies.total_price_cached
CREATE OR REPLACE FUNCTION pc_configurator.recalc_assembly_total_price(
    p_assembly_id BIGINT
)
    RETURNS NUMERIC
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
DECLARE
    v_asm      assemblies%ROWTYPE;
    v_total    NUMERIC := 0;
    v_price    NUMERIC;
    v_ram_id   BIGINT;
    v_drive_id BIGINT;
BEGIN
    -- Загружаем сборку
    SELECT *
    INTO v_asm
    FROM pc_configurator.assemblies
    WHERE assembly_id = p_assembly_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Сборка с id=% не найдена', p_assembly_id
            USING ERRCODE = 'no_data_found';
    END IF;

    ------------------------------------------------------------
    -- 1. Базовые компоненты (одиночные)
    ------------------------------------------------------------
    -- CPU
    v_price := get_min_price('CPU', v_asm.cpu_id);
    IF v_price IS NOT NULL THEN
        v_total := v_total + v_price;
    END IF;

    -- GPU (опционально)
    IF v_asm.gpu_id IS NOT NULL THEN
        v_price := get_min_price('GPU', v_asm.gpu_id);
        IF v_price IS NOT NULL THEN
            v_total := v_total + v_price;
        END IF;
    END IF;

    -- Motherboard
    v_price := get_min_price('MOTHERBOARD', v_asm.motherboard_id);
    IF v_price IS NOT NULL THEN
        v_total := v_total + v_price;
    END IF;

    -- PSU
    v_price := get_min_price('PSU', v_asm.psu_id);
    IF v_price IS NOT NULL THEN
        v_total := v_total + v_price;
    END IF;

    -- Case
    v_price := get_min_price('CASE', v_asm.case_id);
    IF v_price IS NOT NULL THEN
        v_total := v_total + v_price;
    END IF;

    -- Cooler (опционально)
    IF v_asm.cooler_id IS NOT NULL THEN
        v_price := get_min_price('COOLER', v_asm.cooler_id);
        IF v_price IS NOT NULL THEN
            v_total := v_total + v_price;
        END IF;
    END IF;

    ------------------------------------------------------------
    -- 2. RAM (все комплекты сборки)
    ------------------------------------------------------------
    FOR v_ram_id IN
        SELECT ram_kit_id
        FROM pc_configurator.assembly_ram_kits
        WHERE assembly_id = p_assembly_id
        LOOP
            v_price := get_min_price('RAM', v_ram_id);
            IF v_price IS NOT NULL THEN
                v_total := v_total + v_price;
            END IF;
        END LOOP;

    ------------------------------------------------------------
    -- 3. Drives (все накопители сборки)
    ------------------------------------------------------------
    FOR v_drive_id IN
        SELECT drive_id
        FROM pc_configurator.assembly_drives
        WHERE assembly_id = p_assembly_id
        LOOP
            v_price := get_min_price('DRIVE', v_drive_id);
            IF v_price IS NOT NULL THEN
                v_total := v_total + v_price;
            END IF;
        END LOOP;

    ------------------------------------------------------------
    -- 4. Обновляем кэш и возвращаем сумму
    ------------------------------------------------------------
    UPDATE pc_configurator.assemblies
    SET total_price_cached = v_total,
        updated_at         = NOW()
    WHERE assembly_id = p_assembly_id;

    RETURN v_total;
END;
$$;

-- SELECT recalc_assembly_total_price(1);
-- -- вернёт число и обновит assemblies.total_price_cached

------------------------------------------------------------
-- 4) ПРОЦЕДУРЫ-ДОБАВЛЯЛКИ (удобны для API)
------------------------------------------------------------

-- Процедура: добавить RAM-кит в сборку с автопроверкой и пересчётом цены
CREATE OR REPLACE PROCEDURE add_ram_to_assembly(
    p_assembly_id BIGINT,
    p_ram_kit_id BIGINT
)
    LANGUAGE plpgsql
AS
$$
BEGIN
    INSERT INTO assembly_ram_kits (assembly_id, ram_kit_id)
    VALUES (p_assembly_id, p_ram_kit_id);

    -- Если триггер check_assembly_compatibility "ругается",
    -- вставка не завершится и до сюда мы не дойдём.

    PERFORM recalc_assembly_total_price(p_assembly_id);
END;
$$;

-- Процедура: добавить диск в сборку с автопроверкой и пересчётом цены
CREATE OR REPLACE PROCEDURE pc_configurator.add_drive_to_assembly(
    p_assembly_id BIGINT,
    p_drive_id BIGINT,
    p_mount_type VARCHAR DEFAULT NULL -- '2.5', '3.5', 'M.2'
)
    LANGUAGE plpgsql
AS
$$
BEGIN
    INSERT INTO pc_configurator.assembly_drives (assembly_id, drive_id, mount_type)
    VALUES (p_assembly_id, p_drive_id, p_mount_type);

    -- Триггер проверки совместимости сработает автоматически

    PERFORM recalc_assembly_total_price(p_assembly_id);
END;
$$;

------------------------------------------------------------
-- 5) ТРИГГЕРЫ (автоматическая проверка совместимости)
------------------------------------------------------------
-- Триггер: проверка сборки при вставке/обновлении assemblies
CREATE OR REPLACE FUNCTION pc_configurator.trg_assemblies_check()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
BEGIN
    PERFORM pc_configurator.check_assembly_compatibility(NEW.assembly_id);
    RETURN NEW;
END;
$$;

CREATE TRIGGER assemblies_check_compatibility
    AFTER INSERT OR UPDATE
    ON pc_configurator.assemblies
    FOR EACH ROW
EXECUTE FUNCTION pc_configurator.trg_assemblies_check();

-- Триггер для assembly_ram_kits
CREATE OR REPLACE FUNCTION pc_configurator.trg_assembly_ram_kits_check()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
DECLARE
    v_assembly_id BIGINT;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        v_assembly_id := OLD.assembly_id;
    ELSE
        v_assembly_id := NEW.assembly_id;
    END IF;

    PERFORM pc_configurator.check_assembly_compatibility(v_assembly_id);
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$;

CREATE TRIGGER assembly_ram_kits_check_compatibility
    AFTER INSERT OR UPDATE OR DELETE
    ON pc_configurator.assembly_ram_kits
    FOR EACH ROW
EXECUTE FUNCTION pc_configurator.trg_assembly_ram_kits_check();


-- Триггер для assembly_drives
CREATE OR REPLACE FUNCTION pc_configurator.trg_assembly_drives_check()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    SET search_path = pc_configurator
    AS
$$
DECLARE
    v_assembly_id BIGINT;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        v_assembly_id := OLD.assembly_id;
    ELSE
        v_assembly_id := NEW.assembly_id;
    END IF;

    PERFORM check_assembly_compatibility(v_assembly_id);
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$;

CREATE TRIGGER assembly_drives_check_compatibility
    AFTER INSERT OR UPDATE OR DELETE
    ON pc_configurator.assembly_drives
    FOR EACH ROW
EXECUTE FUNCTION pc_configurator.trg_assembly_drives_check();


