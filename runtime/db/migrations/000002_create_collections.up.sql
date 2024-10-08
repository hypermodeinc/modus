CREATE TABLE IF NOT EXISTS "collection_texts" (
    "id" BIGSERIAL PRIMARY KEY,
    "collection" TEXT NOT NULL,
    "key" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS collection_texts_collection_key_idx ON collection_texts (collection, key);

CREATE TABLE IF NOT EXISTS "collection_vectors" (
    "id" BIGSERIAL PRIMARY KEY,
    "text_id" BIGINT REFERENCES collection_texts(id) ON DELETE CASCADE,
    "search_method" TEXT NOT NULL,
    "vector" REAL[] NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS collection_vectors_search_method_text_id_idx ON collection_vectors (search_method, text_id);
