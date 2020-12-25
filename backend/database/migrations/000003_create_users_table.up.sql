CREATE TABLE IF NOT EXISTS "user" 
(
    creation_timestamp timestamp    NOT NULL,
    update_timestamp   timestamp,
    id                 uuid         NOT NULL PRIMARY KEY,
    first_name         VARCHAR      NOT NULL,
    family_name        VARCHAR      NOT NULL,
    age                int          NOT NULL DEFAULT 0 
);
