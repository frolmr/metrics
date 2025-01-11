-- +goose Up
-- +goose StatementBegin
BEGIN;
CREATE TABLE IF NOT EXISTS counter_metrics(
   name VARCHAR (50) UNIQUE NOT NULL,
   value BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS gauge_metrics(
   name VARCHAR (50) UNIQUE NOT NULL,
   value DOUBLE PRECISION NOT NULL
);
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS counter_metrics;
DROP TABLE IF EXISTS gauge_metrics;
-- +goose StatementEnd
