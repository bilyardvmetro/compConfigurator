------------------------------------------------------------
-- 7. ПАРА БАЗОВЫХ ИНДЕКСОВ ПОД ПОИСК
------------------------------------------------------------
SET search_path TO pc_configurator;

CREATE INDEX idx_users_role_id ON users (role_id);
CREATE INDEX idx_assemblies_user_id ON assemblies (user_id);

CREATE INDEX idx_cpus_socket ON cpus (socket_code);
CREATE INDEX idx_motherboards_socket ON motherboards (socket_code);
CREATE INDEX idx_motherboards_form_factor ON motherboards (form_factor_code);

CREATE INDEX idx_ram_kits_type_freq ON ram_kits (ram_type, freq_mhz);
CREATE INDEX idx_storage_drives_type_interface ON storage_drives (drive_type, interface);

CREATE INDEX idx_product_offers_component ON product_offers (component_type, component_id);