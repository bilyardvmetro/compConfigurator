SET search_path TO pc_configurator;

DROP INDEX IF EXISTS idx_product_offers_component;

DROP INDEX IF EXISTS idx_storage_drives_type_interface;
DROP INDEX IF EXISTS idx_ram_kits_type_freq;

DROP INDEX IF EXISTS idx_motherboards_form_factor;
DROP INDEX IF EXISTS idx_motherboards_socket;
DROP INDEX IF EXISTS idx_cpus_socket;

DROP INDEX IF EXISTS idx_assemblies_user_id;
DROP INDEX IF EXISTS idx_users_role_id;
