-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS Limit(
    cid         INTEGER NOT NULL PRIMARY KEY,
    value       INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS Transaction(
    cid         INTEGER NOT NULL,
    tid         TEXT NOT NULL,
    value       INT NOT NULL,
    description TEXT NOT NULL,
    created_at  INTEGER NOT NULL,

    PRIMARY KEY(cid, tid),
    FOREIGN KEY(cid) REFERENCES Limit(cid) ON DELETE CASCADE
);

INSERT INTO Limit (cid, value)
       VALUES (1, 100000), (2, 80000), (3, 1000000), (4, 10000000), (5, 500000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Transaction;
dROP TABLE IF EXISTS Limit;
-- +goose StatementEnd
