-- Тестовые данные для конфигуратора ПК
SET search_path TO pc_configurator;

BEGIN;

------------------------------------------------------------
-- 1. Роли и пользователь
------------------------------------------------------------

INSERT INTO roles (role_id, code, description)
VALUES (1, 'USER', 'Обычный пользователь'),
       (2, 'ADMIN', 'Администратор')
ON CONFLICT (role_id) DO NOTHING;

INSERT INTO users (email, password_hash, nickname, avatar_url, role_id)
VALUES ('test@example.com', 'test-hash', 'TestUser', NULL, 1)
ON CONFLICT (user_id) DO NOTHING;

------------------------------------------------------------
-- 2. Справочники: сокеты и форм-факторы
------------------------------------------------------------

INSERT INTO cpu_sockets (socket_code, description)
VALUES ('AM4', 'Socket AM4 (AMD)'),
       ('LGA1700', 'Socket LGA1700 (Intel)')
ON CONFLICT (socket_code) DO NOTHING;

INSERT INTO motherboard_form_factors (form_factor_code, description, width_mm, height_mm)
VALUES ('ATX', 'Standard ATX', 305, 244),
       ('mATX', 'Micro ATX', 244, 244)
ON CONFLICT (form_factor_code) DO NOTHING;

------------------------------------------------------------
-- 3. Базовые комплектующие
------------------------------------------------------------

-- Кулер
INSERT INTO cpu_coolers (name, brand, cooling_type, tdp_limit_watt, height_mm, fan_count, fan_control_type)
VALUES ('DeepCool AK620', 'DeepCool', 'air', 260, 160, 2, 'PWM')
ON CONFLICT (cooler_id) DO NOTHING;

-- Корпус
INSERT INTO cases (name, brand, psu_form_factor,
                   max_gpu_length_mm, max_gpu_width_slots,
                   max_cooler_height_mm, max_psu_length_mm,
                   drive_bays_3_5_count, drive_bays_2_5_count)
VALUES ('Phanteks P400A', 'Phanteks', 'ATX',
        400, 3.0,
        170, 200,
        2, 2)
ON CONFLICT (case_id) DO NOTHING;

-- Блок питания
INSERT INTO psus (name, brand, power_watt, form_factor,
                  efficiency_rating,
                  pcie_connectors_6_8pin_count,
                  cpu_8pin_connectors_count,
                  sata_connectors_count,
                  molex_connectors_count,
                  length_mm)
VALUES ('Corsair RM650x', 'Corsair', 650, 'ATX',
        '80+ Gold',
        4,
        2,
        8,
        4,
        160)
ON CONFLICT (psu_id) DO NOTHING;

-- Процессор
INSERT INTO cpus (name, brand, socket_code, architecture,
                  core_count, thread_count,
                  base_clock_mhz, boost_clock_mhz,
                  tdp_watt, has_integrated_gpu,
                  supported_ram_type, supported_ram_freq_max_mhz,
                  memory_channels, pcie_version, pcie_lanes_total)
VALUES ('Ryzen 5 5600X', 'AMD', 'AM4', 'Zen 3',
        6, 12,
        3700, 4600,
        65, FALSE,
        'DDR4', 3200,
        2, 4, 20)
ON CONFLICT (cpu_id) DO NOTHING;

-- Видеокарта
INSERT INTO gpus (name, brand, interface,
                  pcie_version, pcie_lanes_required,
                  tdp_watt, power_connectors,
                  length_mm, width_slots, height_mm)
VALUES ('GeForce RTX 3060', 'NVIDIA', 'PCIe x16',
        4, 16,
        170, '1x8-pin',
        242, 2.0, 120)
ON CONFLICT (gpu_id) DO NOTHING;

-- Материнская плата
INSERT INTO motherboards (name, brand,
                          socket_code, form_factor_code,
                          chipset,
                          ram_type, ram_slots,
                          ram_capacity_max_gb, ram_freq_max_mhz,
                          pcie_x16_slots_count, pcie_version_max,
                          m2_slots_count, sata_ports_count,
                          psu_main_connector_type, cpu_power_connector_type)
VALUES ('MSI B550 Tomahawk', 'MSI',
        'AM4', 'ATX',
        'B550',
        'DDR4', 4,
        128, 4400,
        2, 4,
        2, 6,
        '24-pin', '1x8-pin')
ON CONFLICT (motherboard_id) DO NOTHING;

-- ОЗУ
INSERT INTO ram_kits (name, brand,
                      ram_type,
                      module_capacity_gb, module_count, total_capacity_gb,
                      freq_mhz, timings, voltage_v, form_factor)
VALUES ('Kingston Fury 16GB (2x8GB) DDR4-3200', 'Kingston',
        'DDR4',
        8, 2, 16,
        3200, '16-18-18', 1.35, 'DIMM')
ON CONFLICT (ram_kit_id) DO NOTHING;

-- Накопители
INSERT INTO storage_drives (name, brand,
                            drive_type, form_factor, interface, capacity_gb)
VALUES ('Samsung 970 EVO Plus 1TB', 'Samsung',
        'NVMe_SSD', 'M.2', 'NVMe', 1000),
       ('Seagate Barracuda 2TB', 'Seagate',
        'HDD', '3.5', 'SATA', 2000)
ON CONFLICT (drive_id) DO NOTHING;

------------------------------------------------------------
-- 4. Связки по совместимости
------------------------------------------------------------

-- Кулер поддерживает сокет AM4
INSERT INTO cooler_sockets (cooler_id, socket_code, notes)
VALUES (1, 'AM4', '')
ON CONFLICT (cooler_id, socket_code) DO NOTHING;

-- Корпус поддерживает ATX и mATX
INSERT INTO case_form_factor_support (case_id, form_factor_code)
VALUES (1, 'ATX'),
       (1, 'mATX')
ON CONFLICT (case_id, form_factor_code) DO NOTHING;

------------------------------------------------------------
-- 5. Магазины и предложения
------------------------------------------------------------

INSERT INTO shops (name, url)
VALUES ('DNS', 'https://www.dns-shop.ru'),
       ('Citilink', 'https://www.citilink.ru')
ON CONFLICT (shop_id) DO NOTHING;

INSERT INTO product_offers (shop_id, component_type, component_id, price, available)
VALUES (1, 'CPU', 1, 15999, TRUE),
       (2, 'CPU', 1, 14999, TRUE),
       (1, 'GPU', 1, 29999, TRUE),
       (1, 'MOTHERBOARD', 1, 19999, TRUE),
       (1, 'RAM', 1, 5999, TRUE),
       (1, 'DRIVE', 1, 8999, TRUE),
       (1, 'DRIVE', 2, 5999, TRUE),
       (1, 'PSU', 1, 8999, TRUE),
       (1, 'CASE', 1, 6999, TRUE),
       (1, 'COOLER', 1, 3999, TRUE)
ON CONFLICT (offer_id) DO NOTHING;

------------------------------------------------------------
-- 6. Сборка + RAM + диски (совместимая конфигурация)
------------------------------------------------------------

-- Сборка (CPU+GPU+MB+PSU+Case+Cooler)
INSERT INTO assemblies (assembly_id,
                        user_id,
                        name,
                        is_public,
                        total_price_cached,
                        cpu_id,
                        gpu_id,
                        motherboard_id,
                        psu_id,
                        case_id,
                        cooler_id)
VALUES (1,
        2,
        'AM4 Gaming Build',
        TRUE,
        0, -- потом можно обновить суммой из product_offers
        1, -- CPU Ryzen 5 5600X
        1, -- GPU RTX 3060
        1, -- MB B550 Tomahawk
        1, -- PSU RM650x
        1, -- Case P400A
        1 -- Cooler AK620
       )
ON CONFLICT (assembly_id) DO NOTHING;

-- RAM в сборке
INSERT INTO assembly_ram_kits (assembly_id, ram_kit_id)
VALUES (1, 1)
ON CONFLICT (assembly_id, ram_kit_id) DO NOTHING;

-- Диски в сборке
INSERT INTO assembly_drives (assembly_id, drive_id, mount_type)
VALUES (1, 1, 'M.2'), -- NVMe в M.2
       (1, 2, '3.5')  -- HDD в корзину 3.5"
ON CONFLICT (assembly_id, drive_id) DO NOTHING;

COMMIT;
