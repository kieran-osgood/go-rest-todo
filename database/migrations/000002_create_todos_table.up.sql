CREATE TABLE IF NOT EXISTS "todo" 
(
    creation_timestamp timestamp      NOT NULL DEFAULT NOW(),
    update_timestamp   timestamp,
    id                 uuid           NOT NULL PRIMARY KEY,
    text               VARCHAR,
    is_done            boolean        NOT NULL DEFAULT FALSE
);
