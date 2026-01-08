SET search_path TO pc_configurator;

-- триггеры
DROP TRIGGER IF EXISTS assembly_drives_check_compatibility ON assembly_drives;
DROP TRIGGER IF EXISTS assembly_ram_kits_check_compatibility ON assembly_ram_kits;
DROP TRIGGER IF EXISTS assemblies_check_compatibility ON assemblies;

-- триггерные функции
DROP FUNCTION IF EXISTS trg_assembly_drives_check();
DROP FUNCTION IF EXISTS trg_assembly_ram_kits_check();
DROP FUNCTION IF EXISTS trg_assemblies_check();

-- процедуры
DROP PROCEDURE IF EXISTS add_drive_to_assembly(BIGINT, BIGINT, VARCHAR);
DROP PROCEDURE IF EXISTS add_ram_to_assembly(BIGINT, BIGINT);

-- функции цен
DROP FUNCTION IF EXISTS recalc_assembly_total_price(BIGINT);
DROP FUNCTION IF EXISTS get_min_price(VARCHAR, BIGINT);

-- функция мягкой проверки
DROP FUNCTION IF EXISTS is_assembly_compatible(BIGINT);

-- функция строгой проверки
DROP FUNCTION IF EXISTS check_assembly_compatibility(BIGINT);
