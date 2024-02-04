# Altruists

In rare situations, a relay cannot be served from POKT Network. Some sample scenarios for when this can happen:
1. Dispatcher Outages - Gateway server cannot retrieve the necessary information to send a relay
2. Bad overall QoS - Majority of Node operators may have their nodes misconfigured improperly or there is a lack of node operators supporting the chain with respect to load.
3. Chain halts - In extreme conditions, if the chain halts, node operators may stop responding to relay requests

So the Gateway Server will attempt to route the traffic to a backup chain node. This could be any source, for example:
1. Other gateway operators chain urls
2. Centralized Nodes

# Altruist Configuration
Altruist urls are stored in a `chain_id` to `url` mapping inside the `altruists` table.

# Inserting Altruists into the DB
```sql
INSERT INTO altruists (chain_id, url) VALUES ('0021', 'eth-mainnet.nodies.app');
```