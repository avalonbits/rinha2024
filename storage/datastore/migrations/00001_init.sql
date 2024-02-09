-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Transaction(
    cid         INTEGER NOT NULL,
    tid         TEXT NOT NULL,
    value       INT NOT NULL,
    description TEXT NOT NULL,
    created_at  INTEGER NOT NULL,

    PRIMARY KEY(cid, tid)
);

CREATE TABLE IF NOT EXISTS Limit(
    cid         INTEGER NOT NULL PRIMARY KEY,
    value       INTEGER NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Transaction;
-- +goose StatementEnd
