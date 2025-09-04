ALTER INDEX IF EXISTS inferences_model_hash_started_at_idx RENAME TO model_hash_started_at_idx;

DROP INDEX IF EXISTS inferences_plugin_id_idx;
DROP INDEX IF EXISTS inferences_plugin_id_function_idx;
DROP INDEX IF EXISTS inferences_function_idx;

ALTER TABLE IF EXISTS "inferences"
DROP CONSTRAINT inferences_plugins_fkey,
DROP COLUMN "plugin_id",
DROP COLUMN "function";

DROP TABLE IF EXISTS "plugins";
