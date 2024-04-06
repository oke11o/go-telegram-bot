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
