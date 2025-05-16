BEGIN;

CREATE TABLE IF NOT EXISTS "agents" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "data" TEXT,
    "updated" TIMESTAMP(3) WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS agents_name_idx ON agents (name);
CREATE INDEX IF NOT EXISTS agents_status_idx ON agents (status);
CREATE INDEX IF NOT EXISTS agents_updated_idx ON agents (updated);
CREATE INDEX IF NOT EXISTS agents_status_updated_idx ON agents (status, updated);

COMMIT;
