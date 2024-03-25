CREATE TABLE chain_configurations
(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    chain_id VARCHAR NOT NULL UNIQUE,
    pocket_request_timeout_duration VARCHAR NOT NULL,
    altruist_url VARCHAR NOT NULL,
    altruist_request_timeout_duration VARCHAR NOT NULL,
    top_bucket_p90latency_duration VARCHAR NOT NULL,
    height_check_block_tolerance INT NOT NULL,
    data_integrity_check_lookback_height INT NOT NULL
) INHERITS (base_model);

-- Insert an example configuration for Ethereum --
INSERT INTO chain_configurations (chain_id, pocket_request_timeout_duration, altruist_url, altruist_request_timeout_duration, top_bucket_p90latency_duration, height_check_block_tolerance, data_integrity_check_lookback_height) VALUES ('0000', '15s', 'example.com', '30s', '150ms', 100, 25);