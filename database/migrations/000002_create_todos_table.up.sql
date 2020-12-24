CREATE TABLE IF NOT EXISTS "todo" 
(
    creation_timestamp timestamp      NOT NULL,
    update_timestamp   timestamp,
    id                 uuid           NOT NULL PRIMARY KEY,
    text               VARCHAR        NOT NULL,
    is_done            boolean        NOT NULL DEFAULT FALSE
);
