BEGIN;

-- create collections table
CREATE TABLE "collections" (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT UNIQUE NOT NULL
    );

CREATE INDEX collections_names_idx ON collections (name);

-- copy existing collection names
INSERT INTO collections (name)
    SELECT DISTINCT t.collection
    FROM collection_texts t
    LEFT JOIN collections c ON t.collection = c.name
    WHERE c.id is null
    ORDER BY t.collection;

-- create collection_namespaces table
CREATE TABLE "collection_namespaces" (
    "id" SERIAL PRIMARY KEY,
    "collection_id" INTEGER REFERENCES collections(id) ON DELETE CASCADE,
    "name" TEXT NOT NULL
    );

CREATE UNIQUE INDEX collection_namespaces_collection_id_name_idx ON collection_namespaces (collection_id, name);

-- copy existing collection namespaces
INSERT INTO collection_namespaces (collection_id, name)
    SELECT DISTINCT c.id, t.namespace
    FROM collection_texts t
    JOIN collections c ON t.collection = c.name
    LEFT JOIN collection_namespaces n ON c.id = n.collection_id AND t.namespace = n.name
    WHERE n.id is null
    ORDER BY c.id, t.namespace;

-- add namespace_id to collection_texts
ALTER TABLE collection_texts
    ADD COLUMN "namespace_id" INTEGER REFERENCES collection_namespaces(id) ON DELETE CASCADE;

-- populate namespace_id
UPDATE collection_texts t
    SET namespace_id = n.id
    FROM collections c
    JOIN collection_namespaces n ON c.id = n.collection_id
    WHERE t.namespace_id is null AND t.collection = c.name AND t.namespace = n.name;

-- set namespace_id not null
ALTER TABLE collection_texts
    ALTER COLUMN "namespace_id" SET NOT NULL;

-- drop collection and namespace columns
DROP INDEX collection_texts_collection_key_idx;
ALTER TABLE collection_texts
    DROP COLUMN "collection",
    DROP COLUMN "namespace";

-- create the unique index for the namespace_id and key
CREATE UNIQUE INDEX collection_texts_namespace_id_key_idx
    ON collection_texts (namespace_id, key);

COMMIT;