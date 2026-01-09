CREATE TABLE gpus
(
    gpu_id              BIGSERIAL PRIMARY KEY,
    name                VARCHAR(128) NOT NULL,
    brand               VARCHAR(64),
    interface           VARCHAR(32), -- "PCIe x16"
    pcie_version        INTEGER,
    pcie_lanes_required INTEGER,
    tdp_watt            INTEGER,
    power_connectors    VARCHAR(64), -- "2x8-pin"
    length_mm           INTEGER,
    width_slots         NUMERIC(3, 1),
    height_mm           INTEGER
);