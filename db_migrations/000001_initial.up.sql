CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE base_model
(
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE pokt_applications
(
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    encrypted_private_key      BYTEA          NOT NULL
) INHERITS (base_model);



