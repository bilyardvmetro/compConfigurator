CREATE TABLE motherboards
(
    motherboard_id           BIGSERIAL PRIMARY KEY,
    name                     VARCHAR(128) NOT NULL,
    brand                    VARCHAR(64),
    socket_code              VARCHAR(32)  NOT NULL REFERENCES cpu_sockets (socket_code),
    form_factor_code         VARCHAR(16)  NOT NULL REFERENCES motherboard_form_factors (form_factor_code),
    chipset                  VARCHAR(64),
    ram_type                 VARCHAR(16), -- DDR4 / DDR5
    ram_slots                INTEGER,
    ram_capacity_max_gb      INTEGER,
    ram_freq_max_mhz         INTEGER,
    pcie_x16_slots_count     INTEGER,
    pcie_version_max         INTEGER,
    m2_slots_count           INTEGER,
    sata_ports_count         INTEGER,
    psu_main_connector_type  VARCHAR(32), -- "24-pin"
    cpu_power_connector_type VARCHAR(32)  -- "1x8-pin", "2x8-pin"
);