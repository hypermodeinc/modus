CREATE TABLE IF NOT EXISTS "local_instance" (
    model_name TEXT NOT NULL,
    model_task TEXT NOT NULL,
    source_model TEXT NOT NULL,
    model_provider TEXT NOT NULL,
    model_host TEXT NOT NULL,
    model_version TEXT NOT NULL,
    model_hash TEXT NOT NULL,
    input TEXT NOT NULL,
    output TEXT NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);