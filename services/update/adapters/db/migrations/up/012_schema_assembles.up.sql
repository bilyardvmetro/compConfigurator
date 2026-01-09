CREATE TABLE assemblies
(
    assembly_id        BIGSERIAL PRIMARY KEY,
    user_id            BIGINT       NOT NULL REFERENCES users (user_id),
    name               VARCHAR(128) NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP    NOT NULL DEFAULT NOW(),
    is_public          BOOLEAN      NOT NULL DEFAULT FALSE,
    total_price_cached NUMERIC(12, 2),

    cpu_id             BIGINT       NOT NULL REFERENCES cpus (cpu_id),
    gpu_id             BIGINT REFERENCES gpus (gpu_id),
    motherboard_id     BIGINT       NOT NULL REFERENCES motherboards (motherboard_id),
    psu_id             BIGINT       NOT NULL REFERENCES psus (psu_id),
    case_id            BIGINT       NOT NULL REFERENCES cases (case_id),
    cooler_id          BIGINT REFERENCES cpu_coolers (cooler_id)
);