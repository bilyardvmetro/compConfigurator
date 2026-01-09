CREATE TABLE cases
(
    case_id              BIGSERIAL PRIMARY KEY,
    name                 VARCHAR(128) NOT NULL,
    brand                VARCHAR(64),
    psu_form_factor      VARCHAR(16), -- ATX / SFX
    max_gpu_length_mm    INTEGER,
    max_gpu_width_slots  NUMERIC(3, 1),
    max_cooler_height_mm INTEGER,
    max_psu_length_mm    INTEGER,
    drive_bays_3_5_count INTEGER,
    drive_bays_2_5_count INTEGER
);