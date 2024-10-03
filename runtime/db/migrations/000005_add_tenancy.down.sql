BEGIN;

ALTER TABLE collection_texts DROP COLUMN IF EXISTS namespace;

COMMIT;
