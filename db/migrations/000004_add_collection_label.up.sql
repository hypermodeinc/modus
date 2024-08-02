BEGIN;

ALTER TABLE collection_texts ADD COLUMN label TEXT;

COMMIT;