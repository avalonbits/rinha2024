-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Transaction(
    cid         INTEGER NOT NULL,
    tid         TEXT NOT NULL,
    value       INTEGER NOT NULL,
    description TEXT NOT NULL,
    created_at  INTEGER NOT NULL,

    PRIMARY KEY(cid, tid)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Transaction;
-- +goose StatementEnd
