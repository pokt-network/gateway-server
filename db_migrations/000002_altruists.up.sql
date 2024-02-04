CREATE TABLE altruists
(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    chain_id VARCHAR NOT NULL UNIQUE,
    url VARCHAR NOT NULL
) INHERITS (base_model);