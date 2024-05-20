CREATE TABLE IF NOT EXISTS "local_instance" (
    model_id UUID NOT NULL,
    input TEXT NOT NULL,
    output TEXT NOT NULL,
    started_at TIMESTAMP(3) WITH TIME ZONE NOT NULL,
    duration_ms INTEGER NOT NULL
);
