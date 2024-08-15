# POKT Gateway Server Architecture <!-- omit in toc -->

![gateway-server-architecture.png](resources/gateway-server-architecture.png)

- [Gateway Server Responsibilities](#gateway-server-responsibilities)
  - [Primary Features](#primary-features)
  - [Secondary Features](#secondary-features)
- [Gateway Operator Responsibilities](#gateway-operator-responsibilities)

## Gateway Server Responsibilities

### Primary Features

Under the hood, the gateway server handles everything in regard to protocol interaction to abstract away the complexity of:

1. Retrieving a session
2. Signing a relay
3. Sending a relay to a node operator & receiving a response

### Secondary Features

1. **Node Selection & Routing (QoS)** - determining which nodes are healthy based off responses
2. **Metrics** - Provides underlying Prometheus metrics endpoint for relay performance metadata
3. **HTTP Interface** - Providing an efficient HTTP endpoint to send requests to

## Gateway Operator Responsibilities

1. **Key management** - Keeping the encryption key and respectively the app stakes keys secure.
2. **App stake management** - Staking in the approriate chains
3. **SaaS business support** - Any features in regard to a SaaS business as mentioned in the [overview](overview.md).
