CREATE TABLE IF NOT EXISTS "plugins" (
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
    "version" TEXT,
    "language" TEXT NOT NULL,
    "sdk_version" TEXT NOT NULL,
    "build_id" TEXT UNIQUE NOT NULL,
    "build_time" TIMESTAMP(3) WITH TIME ZONE NOT NULL,
    "git_repo" TEXT,
    "git_commit" TEXT
);

CREATE INDEX IF NOT EXISTS plugins_build_id_idx ON plugins (build_id);
CREATE INDEX IF NOT EXISTS plugins_build_time_idx ON plugins (build_time);

ALTER TABLE IF EXISTS "inferences"
ADD COLUMN "plugin_id" UUID,
ADD COLUMN "function" TEXT,
ADD CONSTRAINT inferences_plugins_fkey FOREIGN KEY (plugin_id) REFERENCES plugins (id);

CREATE INDEX IF NOT EXISTS inferences_plugin_id_idx ON inferences (plugin_id);
CREATE INDEX IF NOT EXISTS inferences_plugin_id_function_idx ON inferences (plugin_id, function);
CREATE INDEX IF NOT EXISTS inferences_function_idx ON inferences (function);

ALTER INDEX IF EXISTS model_hash_started_at_idx RENAME TO inferences_model_hash_started_at_idx;
