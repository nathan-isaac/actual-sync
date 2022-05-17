CREATE TABLE IF NOT EXISTS messages_binary
(
    timestamp TEXT PRIMARY KEY,
    is_encrypted BOOLEAN,
    content BLOB
);

CREATE TABLE IF NOT EXISTS messages_merkles
(
    id TEXT PRIMARY KEY,
    merkle TEXT
);