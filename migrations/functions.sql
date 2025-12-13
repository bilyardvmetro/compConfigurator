SET search_path TO pc_configurator;

-- Функция: вернуть признак совместимости и текст ошибки (если есть)
CREATE OR REPLACE FUNCTION is_assembly_compatible(
    p_assembly_id BIGINT,
    OUT is_ok BOOLEAN,
    OUT error_message TEXT
)
    RETURNS RECORD AS
$$
BEGIN
    -- Пытаемся вызвать "строгую" проверку, которая кидает исключения
    PERFORM check_assembly_compatibility(p_assembly_id);

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
$$ LANGUAGE plpgsql;

-- SELECT * FROM is_assembly_compatible(1);
-- -- вернёт, например:  is_ok = true, error_message = null

-- Вспомогательная функция: минимальная цена для конкретного компонента
CREATE OR REPLACE FUNCTION get_min_price(
    p_component_type VARCHAR,
    p_component_id BIGINT
)
    RETURNS NUMERIC AS
$$
DECLARE
    v_price NUMERIC;
BEGIN
    SELECT MIN(price)
    INTO v_price
    FROM product_offers
    WHERE component_type = p_component_type
      AND component_id = p_component_id
      AND available = TRUE;

    RETURN v_price; -- может быть NULL, если предложений нет
END;
$$ LANGUAGE plpgsql;

-- Пересчитать общую цену сборки и записать в assemblies.total_price_cached
CREATE OR REPLACE FUNCTION recalc_assembly_total_price(
    p_assembly_id BIGINT
)
    RETURNS NUMERIC AS
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
    FROM assemblies
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
        FROM assembly_ram_kits
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
        FROM assembly_drives
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
    UPDATE assemblies
    SET total_price_cached = v_total,
        updated_at         = NOW()
    WHERE assembly_id = p_assembly_id;

    RETURN v_total;
END;
$$ LANGUAGE plpgsql;

-- SELECT recalc_assembly_total_price(1);
-- -- вернёт число и обновит assemblies.total_price_cached

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
CREATE OR REPLACE PROCEDURE add_drive_to_assembly(
    p_assembly_id BIGINT,
    p_drive_id BIGINT,
    p_mount_type VARCHAR DEFAULT NULL -- '2.5', '3.5', 'M.2'
)
    LANGUAGE plpgsql
AS
$$
BEGIN
    INSERT INTO assembly_drives (assembly_id, drive_id, mount_type)
    VALUES (p_assembly_id, p_drive_id, p_mount_type);

    -- Триггер проверки совместимости сработает автоматически

    PERFORM recalc_assembly_total_price(p_assembly_id);
END;
$$;

