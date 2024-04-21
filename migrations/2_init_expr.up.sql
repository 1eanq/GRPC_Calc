CREATE TABLE IF NOT EXISTS expressions
(
    id        INTEGER PRIMARY KEY,
    expression     TEXT NOT NULL UNIQUE,
    answer TEXT,
    uid INTEGER
);

