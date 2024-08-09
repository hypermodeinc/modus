BEGIN;

-- restore columns
ALTER TABLE collection_texts
    ADD COLUMN "collection" TEXT,
    ADD COLUMN "namespace" TEXT DEFAULT '' NOT NULL;

CREATE INDEX collection_texts_collection_key_idx ON collection_texts (collection, key);

-- re-populate collection and namespace columns
UPDATE collection_texts t
    SET collection = c.name, namespace = n.name
    FROM collections c
    JOIN collection_namespaces n ON c.id = n.collection_id
    WHERE t.namespace_id = n.id;

ALTER TABLE collection_texts ALTER COLUMN "collection" SET NOT NULL;


-- drop new column and tables
ALTER TABLE collection_texts DROP COLUMN "namespace_id";
DROP TABLE collection_namespaces;
DROP TABLE collections;

-- drop helper function
DROP FUNCTION unnest_nd_1d;

COMMIT;
