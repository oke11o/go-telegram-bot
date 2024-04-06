CREATE TABLE IF NOT EXISTS income_request
(
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    from_id             INTEGER NOT NULL,
    message_id          INTEGER NOT NULL,
    reply_to_message_id INTEGER NOT NULL,
    request_id          VARCHAR NOT NULL,
    username            VARCHAR NOT NULL,
    message             TEXT    NOT NULL,
    text                VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS user
(
    id            INTEGER PRIMARY KEY,
    username      VARCHAR NOT NULL,
    first_name    VARCHAR NOT NULL,
    last_name     VARCHAR NOT NULL,
    language_code VARCHAR NOT NULL,
    is_bot        INTEGER NOT NULL,
    is_maintainer INTEGER NOT NULL,
    is_manager    INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS chat
(
    id       INTEGER PRIMARY KEY,
    type     VARCHAR,
    title    VARCHAR,
    photo    VARCHAR,
    location VARCHAR
);

CREATE TABLE IF NOT EXISTS tournament
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      VARCHAR NOT NULL,
    date       VARCHAR NOT NULL,
    status     VARCHAR NOT NULL,
    created_by INTEGER NOT NULL,
    created_at VARCHAR NOT NULL,
    updated_at VARCHAR NOT NULL,
    FOREIGN KEY (created_by) REFERENCES user (id)
);

CREATE TABLE IF NOT EXISTS session
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    status     VARCHAR NOT NULL,
    data       TEXT    NOT NULL,
    created_at VARCHAR NOT NULL,
    updated_at VARCHAR NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (id)

);
