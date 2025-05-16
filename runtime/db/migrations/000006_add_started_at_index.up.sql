BEGIN;

CREATE INDEX idx_started_at_desc ON inferences (started_at DESC);

COMMIT;
