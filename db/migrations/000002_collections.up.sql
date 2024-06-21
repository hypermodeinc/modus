CREATE TABLE IF NOT EXISTS "collection_texts" (
    "id" BIGSERIAL PRIMARY KEY,
    "collection" TEXT NOT NULL,
    "key" TEXT NOT NULL,
    "text" TEXT NOT NULL
    );

CREATE INDEX IF NOT EXISTS collection_texts_id_collection_idx ON collection_texts (id, collection);

CREATE TABLE IF NOT EXISTS "collection_vectors" (
    "id" BIGSERIAL PRIMARY KEY,
    "text_id" BIGINT REFERENCES collection_texts(id) ON DELETE CASCADE,
    "search_method" TEXT NOT NULL,
    "vector" REAL[] NOT NULL
    );

CREATE INDEX IF NOT EXISTS collection_vectors_id_search_method_text_id_idx ON collection_vectors (id, search_method, text_id);