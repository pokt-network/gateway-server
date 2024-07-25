# POKT Gateway Server API Endpoints <!-- omit in toc -->

The Gateway Server currently exposes all its API endpoints in form of HTTP endpoints.

Postman collection can be found [here](https://www.postman.com/dark-shadow-851601/workspace/os-gateway/collection/27302708-537f3ba3-3193-4290-98d0-0d5836988a2f)

`x-api-key` is an api key set by the gateway operator to transmit internal private data

_TODO_IMPROVE: Move this to Swagger in the future if our API endpoints become more complex._

- [API Endpoints](#api-endpoints)
- [Examples](#examples)
  - [Relay](#relay)
  - [Metrics](#metrics)
  - [PoktApps](#poktapps)
    - [List](#list)
    - [Add](#add)
    - [Delete](#delete)
    - [QoS Noes](#qos-noes)

## API Endpoints

| Endpoint             | HTTP METHOD | Description                                                                                                                                      | HEADERS     | Request Parameters                       |
| -------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ----------- | ---------------------------------------- |
| `/relay/{chain_id}`  | ANY         | The main endpoint to send relays to                                                                                                              | ANY         | `{chain_id}` - Network identifier        |
| `/metrics`           | GET         | Gateway metadata related to server performance and observability                                                                                 | N/A         | N/A                                      |
| `/poktapps`          | GET         | List all the available app stakes                                                                                                                | `x-api-key` | N/A                                      |
| `/poktapps`          | POST        | Add an existing app stake to the appstake database (not recommended due to security)                                                             | `x-api-key` | `private_key` - private key of app stake |
| `/poktapps/{app_id}` | DELETE      | Remove an existing app stake from the appstake database (not recommended due to security)                                                        | `x-api-key` | `app_id` - id of the appstake            |
| `/qosnodes`          | GET         | List of nodes and public QoS state such as healthiness and last known error. This can be used to expose to node operators to improve visibility. | `x-api-key` | N/A                                      |

## Examples

These examples assume gateway server is running locally.

### Relay

```bash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8080/relay/0021
```

### Metrics

```bash
curl -X GET http://localhost:8080/metrics
```

### PoktApps

Make sure that the gateway server starts up with the `API_KEY` environment variable set.

#### List

```bash
curl -X GET http://localhost:8080/poktapps
```

#### Add

```bash
curl -X POST -H "x-api-key: $API_KEY" https://localhost:8080/poktapps/{private_key}
```

#### Delete

```bash
curl -X DELETE -H "x-api-key: $API_KEY" https://localhost:8080/poktapps/{app_id}
```

#### QoS Noes

```bash
curl -X GET -H "x-api-key: $API_KEY" http://localhost:8080/qosnodes
```
