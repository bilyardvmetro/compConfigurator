CREATE TABLE ram_kits
(
    ram_kit_id         BIGSERIAL PRIMARY KEY,
    name               VARCHAR(128) NOT NULL,
    brand              VARCHAR(64),
    ram_type           VARCHAR(16), -- DDR4 / DDR5
    module_capacity_gb INTEGER,
    module_count       INTEGER,
    total_capacity_gb  INTEGER,
    freq_mhz           INTEGER,
    timings            VARCHAR(32),
    voltage_v          NUMERIC(3, 2),
    form_factor        VARCHAR(16)  -- DIMM / SO-DIMM
);