BEGIN;

ALTER TABLE collection_texts ADD COLUMN labels TEXT[];

COMMIT;
