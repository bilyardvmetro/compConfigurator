-- Триггер: проверка сборки при вставке/обновлении assemblies
CREATE OR REPLACE FUNCTION trg_assemblies_check()
    RETURNS TRIGGER AS
$$
BEGIN
    PERFORM check_assembly_compatibility(NEW.assembly_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER assemblies_check_compatibility
    AFTER INSERT OR UPDATE
    ON assemblies
    FOR EACH ROW
EXECUTE FUNCTION trg_assemblies_check();

-- Триггер для assembly_ram_kits
CREATE OR REPLACE FUNCTION trg_assembly_ram_kits_check()
    RETURNS TRIGGER AS
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
$$ LANGUAGE plpgsql;

CREATE TRIGGER assembly_ram_kits_check_compatibility
    AFTER INSERT OR UPDATE OR DELETE
    ON assembly_ram_kits
    FOR EACH ROW
EXECUTE FUNCTION trg_assembly_ram_kits_check();


-- Триггер для assembly_drives
CREATE OR REPLACE FUNCTION trg_assembly_drives_check()
    RETURNS TRIGGER AS
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
$$ LANGUAGE plpgsql;

CREATE TRIGGER assembly_drives_check_compatibility
    AFTER INSERT OR UPDATE OR DELETE
    ON assembly_drives
    FOR EACH ROW
EXECUTE FUNCTION trg_assembly_drives_check();

