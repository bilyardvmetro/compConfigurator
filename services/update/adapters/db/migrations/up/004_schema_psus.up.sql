CREATE TABLE psus
(
    psu_id                       BIGSERIAL PRIMARY KEY,
    name                         VARCHAR(128) NOT NULL,
    brand                        VARCHAR(64),
    power_watt                   INTEGER,
    form_factor                  VARCHAR(16), -- ATX / SFX
    efficiency_rating            VARCHAR(16), -- 80+ Bronze / Gold / ...
    pcie_connectors_6_8pin_count INTEGER,
    cpu_8pin_connectors_count    INTEGER,
    sata_connectors_count        INTEGER,
    molex_connectors_count       INTEGER,
    length_mm                    INTEGER
);