# POKT Gateway Server API Endpoints

The Gateway Server currently exposes all its API endpoints in form of HTTP endpoints. Note, we can move this to Swagger in the future if our API endpoints become more complex.

`x-api-key`  is an api key set by the gateway operator to transmit internal private data

| Endpoint                    | HTTP METHOD | Description                                                           | HEADERS     | Request Parameters                                                  |
|-----------------------------|-------------|-----------------------------------------------------------------------|-------------|---------------------------------------------------------------------|
| `/relay/{chain_id}`         | ANY         | The main endpoint for your reverse proxy to send requests too         | ANY         | `{chain_id}` - Network identifier                                   |
| `/metrics`                  | GET         | Metadata on the gateway server performance for observability purposes | N/A         | N/A                                                                 |
| `/poktapp`                  | POST        | Adds an existing app stake to the appstake database                   | `x-api-key` | `encrypted_private_key` - private key of app stake that's encrypted |
| `/poktapp`                  | DELETE      | REmoves an existing app stake from the appstake database              | `x-api-key` | `address` - address of the appstake that's encrypted                |
| `/poktapp`                  | GET         | A list of all the avaliable app stakes                                | `x-api-key` | N/A                                                                 |
| `/poktapp/{app_id}`         | GET         | Return a specific app stake                                           | `x-api-key` | `app_id` - internal app id                                          |
| `/appstakes/{app_id}/chains` | POST        | Edits the app stake chains                                            | `x-api-key` | `chain_ids` - an string array of chain ids                          |
