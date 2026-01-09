-- Кулер ↔ сокет
CREATE TABLE cooler_sockets
(
    cooler_id   BIGINT      NOT NULL REFERENCES cpu_coolers (cooler_id) ON DELETE CASCADE,
    socket_code VARCHAR(32) NOT NULL REFERENCES cpu_sockets (socket_code) ON DELETE CASCADE,
    notes       TEXT,
    PRIMARY KEY (cooler_id, socket_code)
);

-- Корпус ↔ форм-фактор платы
CREATE TABLE case_form_factor_support
(
    case_id          BIGINT      NOT NULL REFERENCES cases (case_id) ON DELETE CASCADE,
    form_factor_code VARCHAR(16) NOT NULL REFERENCES motherboard_form_factors (form_factor_code) ON DELETE CASCADE,
    PRIMARY KEY (case_id, form_factor_code)
);

-- Сборка ↔ RAM
CREATE TABLE assembly_ram_kits
(
    assembly_id BIGINT NOT NULL REFERENCES assemblies (assembly_id) ON DELETE CASCADE,
    ram_kit_id  BIGINT NOT NULL REFERENCES ram_kits (ram_kit_id),
    PRIMARY KEY (assembly_id, ram_kit_id)
);

-- Сборка ↔ диски
CREATE TABLE assembly_drives
(
    assembly_id BIGINT NOT NULL REFERENCES assemblies (assembly_id) ON DELETE CASCADE,
    drive_id    BIGINT NOT NULL REFERENCES storage_drives (drive_id),
    mount_type  VARCHAR(16), -- "2.5", "3.5", "M.2"
    PRIMARY KEY (assembly_id, drive_id)
);
