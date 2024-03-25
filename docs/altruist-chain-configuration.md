# Altruist (Failover) Request
In rare situations, a relay cannot be served from POKT Network. Some sample scenarios for when this can happen:
1. Dispatcher Outages - Gateway server cannot retrieve the necessary information to send a relay
2. Bad overall QoS - Majority of Node operators may have their nodes misconfigured improperly or there is a lack of node operators supporting the chain with respect to load.
3. Chain halts - In extreme conditions, if the chain halts, node operators may stop responding to relay requests

So the Gateway Server will attempt to route the traffic to a backup chain node. This could be any source, for example:
1. Other gateway operators chain urls
2. Centralized Nodes

# Chain Configuration
Given that every chain has differences and have different sources for fail over, the gateway server allows for optional customization for request timeouts, failover relay, and QoS checks.
The data is stored inside the `chain_configuration` table and is accessed via the [chain_configurations_registry_service.go](..%2Finternal%2Fchain_configurations_registry%2Fchain_configurations_registry_service.go).

_In the event that a config is not provided, the gateway server will assume defaults provided from the specified [config_provider.go](..%2Finternal%2Fglobal_config%2Fconfig_provider.go) and the provided QoS [checks](..%2Finternal%2Fnode_selector_service%2Fchecks)_

# Inserting a custom chain configuration
```sql
-- Insert an example configuration for Ethereum --
INSERT INTO chain_configurations (chain_id, pocket_request_timeout_duration, altruist_url, altruist_request_timeout_duration, top_bucket_p90latency_duration, height_check_block_tolerance, data_integrity_check_lookback_height) VALUES ('0000', '15s', 'example.com', '30s', '150ms', 100, 25);
```

- `chain_id` - id of the Pocket Network Chain
- `pocket_request_time` - duration of the maximum amount of time for a network relay to respond
- `altruist_url` -  source of the relay in the event that a network request fails
- `altruist_request_timeout_duration` - duration of the maximum amount of time for a backup request to respond
- `top_bucket_p90latency_duration` -  maximum amount of latency for nodes to be favored 0 <= x <= `top_bucket_p90latency_duration`
- `height_check_block_tolerance` - number of blocks a node is allowed to be behind (some chains may have node operators moving faster than others)
- `data_integrity_check_lookback_height` - number of blocks data integrity will look behind for source of truth block for other node operators to attest too
