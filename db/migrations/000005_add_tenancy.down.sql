BEGIN;

-- Step 1: Drop the foreign key constraint from collection_vectors
ALTER TABLE collection_vectors DROP CONSTRAINT IF EXISTS collection_vectors_text_id_fkey;

-- Step 2: Drop the namespace column from collection_vectors
ALTER TABLE collection_vectors DROP COLUMN IF EXISTS namespace;

-- Step 3: Drop the composite primary key in collection_texts
ALTER TABLE collection_texts DROP CONSTRAINT IF EXISTS collection_texts_pk;

-- Step 4: Re-add the old primary key constraint
ALTER TABLE collection_texts ADD CONSTRAINT collection_texts_pkey PRIMARY KEY (id);

-- Step 5: Drop the namespace column from collection_texts
ALTER TABLE collection_texts DROP COLUMN IF EXISTS namespace;

-- Step 6: Re-add the foreign key constraint to collection_vectors referencing the id in collection_texts
ALTER TABLE collection_vectors ADD CONSTRAINT collection_vectors_text_id_fkey
    FOREIGN KEY (text_id) REFERENCES collection_texts(id) ON DELETE CASCADE;

COMMIT;
