-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS Limits(
    cid         INTEGER NOT NULL PRIMARY KEY,
    value       INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS Transactions(
    cid         INTEGER NOT NULL,
    tid         TEXT NOT NULL,
    value       INTEGER NOT NULL,
    description TEXT NOT NULL,
    created_at  INTEGER NOT NULL,

    PRIMARY KEY(cid, tid),
    FOREIGN KEY(cid) REFERENCES Limits(cid) ON DELETE CASCADE
);

INSERT INTO Limits (cid, value)
       VALUES (1, 100000), (2, 80000), (3, 1000000), (4, 10000000), (5, 500000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Transaction;
dROP TABLE IF EXISTS Limits;
-- +goose StatementEnd
