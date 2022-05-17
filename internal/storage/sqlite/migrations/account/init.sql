CREATE TABLE IF NOT EXISTS auth
(
    password TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS sessions
(
    token TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS files
(
    id TEXT PRIMARY KEY,
    group_id TEXT,
    sync_version SMALLINT,
    encrypt_meta TEXT,
    encrypt_keyid TEXT,
    encrypt_salt TEXT,
    encrypt_test TEXT,
    deleted BOOLEAN DEFAULT FALSE,
    name TEXT
);