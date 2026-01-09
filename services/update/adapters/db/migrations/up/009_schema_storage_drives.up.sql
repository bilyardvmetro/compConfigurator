CREATE TABLE storage_drives
(
    drive_id    BIGSERIAL PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    brand       VARCHAR(64),
    drive_type  VARCHAR(16), -- HDD / SATA_SSD / NVMe_SSD
    form_factor VARCHAR(16), -- "2.5", "3.5", "M.2"
    interface   VARCHAR(16), -- SATA / NVMe
    capacity_gb INTEGER
);