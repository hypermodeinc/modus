CREATE TABLE IF NOT EXISTS "inferences" (
    "id" UUID PRIMARY KEY,
    "model_hash" TEXT NOT NULL,
    "input" JSONB NOT NULL,
    "output" JSONB NOT NULL,
    "started_at" TIMESTAMP(3) WITH TIME ZONE NOT NULL,
    "duration_ms" INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS model_hash_started_at_idx ON inferences (model_hash, started_at);
