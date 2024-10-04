CREATE TABLE IF NOT EXISTS users
(
    "id"
        INTEGER PRIMARY KEY
        AUTOINCREMENT
        NOT NULL,
    "user"
        TEXT
        NOT NULL,
    "password"
        TEXT
        NOT NULL,
    "active"
        BOOLEAN,
    "requestID"
        TEXT
        NOT NULL
);

CREATE TABLE IF NOT EXISTS reqFormIP
(
    "ip"
        TEXT
        NOT NULL
        PRIMARY KEY,
    "reqTime"
        INTEGER
        NOT NULL
);

CREATE TABLE IF NOT EXISTS requests
(
    "id"
        TEXT
        NOT NULL
        PRIMARY KEY,
    "company"
        TEXT
        NOT NULL,
    "email"
        TEXT
        NOT NULL,
    "applicant"
        TEXT
        NOT NULL,
    "active"
        BOOLEAN
        NOT NULL,
    "requestTime"
        INTEGER
        NOT NULL,
    "submitTime"
        INTEGER
);

PRAGMA journal_mode= WAL;