CREATE TABLE IF NOT EXISTS users(
    user_id INTEGER PRIMARY KEY,
    login TEXT UNIQUE,
    password BLOB,
    is_admin INT
);