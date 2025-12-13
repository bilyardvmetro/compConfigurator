-- Опционально: отдельная схема
CREATE SCHEMA IF NOT EXISTS pc_configurator;
SET search_path TO pc_configurator;

------------------------------------------------------------
-- 1. РОЛИ И ПОЛЬЗОВАТЕЛИ
------------------------------------------------------------

CREATE TABLE roles
(
    role_id     BIGSERIAL PRIMARY KEY,
    code        VARCHAR(32) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE users
(
    user_id       BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname      VARCHAR(64),
    avatar_url    VARCHAR(255),
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    role_id       BIGINT       NOT NULL REFERENCES roles (role_id)
);

------------------------------------------------------------
-- 2. СПРАВОЧНИКИ ДЛЯ СОВМЕСТИМОСТИ
------------------------------------------------------------

CREATE TABLE cpu_sockets
(
    socket_code VARCHAR(32) PRIMARY KEY,
    description TEXT
);

CREATE TABLE motherboard_form_factors
(
    form_factor_code VARCHAR(16) PRIMARY KEY, -- ATX, mATX, ITX
    description      TEXT,
    width_mm         INTEGER,
    height_mm        INTEGER
);

------------------------------------------------------------
-- 3. БАЗОВЫЕ КОМПЛЕКТУЮЩИЕ
------------------------------------------------------------

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

CREATE TABLE cpus
(
    cpu_id                     BIGSERIAL PRIMARY KEY,
    name                       VARCHAR(128) NOT NULL,
    brand                      VARCHAR(64),
    socket_code                VARCHAR(32)  NOT NULL REFERENCES cpu_sockets (socket_code),
    architecture               VARCHAR(64),
    core_count                 INTEGER,
    thread_count               INTEGER,
    base_clock_mhz             INTEGER,
    boost_clock_mhz            INTEGER,
    tdp_watt                   INTEGER,
    has_integrated_gpu         BOOLEAN,
    supported_ram_type         VARCHAR(16), -- DDR4 / DDR5
    supported_ram_freq_max_mhz INTEGER,
    memory_channels            INTEGER,
    pcie_version               INTEGER,
    pcie_lanes_total           INTEGER
);

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

------------------------------------------------------------
-- 4. СБОРКИ
------------------------------------------------------------

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

------------------------------------------------------------
-- 5. СВЯЗУЮЩИЕ M:N (СОВМЕСТИМОСТЬ)
------------------------------------------------------------

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

------------------------------------------------------------
-- 6. МАГАЗИНЫ И ПРЕДЛОЖЕНИЯ
------------------------------------------------------------

CREATE TABLE shops
(
    shop_id BIGSERIAL PRIMARY KEY,
    name    VARCHAR(128) NOT NULL,
    url     VARCHAR(255)
);

CREATE TABLE product_offers
(
    offer_id       BIGSERIAL PRIMARY KEY,
    shop_id        BIGINT         NOT NULL REFERENCES shops (shop_id) ON DELETE CASCADE,
    component_type VARCHAR(32)    NOT NULL, -- 'CPU','GPU','MOTHERBOARD',...
    component_id   BIGINT         NOT NULL,
    price          NUMERIC(12, 2) NOT NULL,
    available      BOOLEAN        NOT NULL DEFAULT TRUE
    -- Полиморфный component_id.
    -- Ссылки на конкретные таблицы будем контролировать на уровне бизнес-логики/триггеров.
);

------------------------------------------------------------
-- 7. ПАРА БАЗОВЫХ ИНДЕКСОВ ПОД ПОИСК
------------------------------------------------------------

CREATE INDEX idx_users_role_id ON users (role_id);
CREATE INDEX idx_assemblies_user_id ON assemblies (user_id);

CREATE INDEX idx_cpus_socket ON cpus (socket_code);
CREATE INDEX idx_motherboards_socket ON motherboards (socket_code);
CREATE INDEX idx_motherboards_form_factor ON motherboards (form_factor_code);

CREATE INDEX idx_ram_kits_type_freq ON ram_kits (ram_type, freq_mhz);
CREATE INDEX idx_storage_drives_type_interface ON storage_drives (drive_type, interface);

CREATE INDEX idx_product_offers_component ON product_offers (component_type, component_id);
