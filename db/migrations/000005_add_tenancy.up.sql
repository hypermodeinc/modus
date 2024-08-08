BEGIN;

-- Step 1: Temporarily drop the foreign key constraint from collection_vectors
ALTER TABLE collection_vectors DROP CONSTRAINT collection_vectors_text_id_fkey;

-- Step 2: Add the namespace column with a default value and set it to NOT NULL in collection_texts
ALTER TABLE collection_texts ADD COLUMN namespace TEXT DEFAULT '' NOT NULL;

-- Update the namespace column for existing rows to the empty string
UPDATE collection_texts SET namespace = '' WHERE namespace IS NULL;

-- Step 3: Drop the existing primary key constraint
ALTER TABLE collection_texts DROP CONSTRAINT collection_texts_pkey;

-- Step 4: Add the new composite primary key
ALTER TABLE collection_texts ADD CONSTRAINT collection_texts_pk PRIMARY KEY (namespace, id);

-- Step 5: Add the namespace column to collection_vectors with a default value
ALTER TABLE collection_vectors ADD COLUMN namespace TEXT DEFAULT '' NOT NULL;

-- Update the namespace column for existing rows to the empty string
UPDATE collection_vectors SET namespace = (SELECT namespace FROM collection_texts WHERE collection_texts.id = collection_vectors.text_id);

-- Step 6: Add the foreign key constraint to collection_vectors referencing the composite primary key in collection_texts
ALTER TABLE collection_vectors ADD CONSTRAINT collection_vectors_text_id_fkey
    FOREIGN KEY (namespace, text_id) REFERENCES collection_texts(namespace, id) ON DELETE CASCADE;

COMMIT;
