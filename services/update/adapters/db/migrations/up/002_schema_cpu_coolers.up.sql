CREATE TABLE cpu_coolers
(
    cooler_id        BIGSERIAL PRIMARY KEY,
    name             VARCHAR(128) NOT NULL,
    brand            VARCHAR(64),
    cooling_type     VARCHAR(16), -- air / liquid
    tdp_limit_watt   INTEGER,
    height_mm        INTEGER,
    fan_count        INTEGER,
    fan_control_type VARCHAR(16)  -- PWM / DC
);