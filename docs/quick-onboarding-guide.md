# Mainnet Onboarding Guide

## 1. Retrieve App Stakes Private Keys
Retrieve your app stakes private keys from the Pocket Network Foundation (PNF).

## 2. Create Encryption Password
Create an encryption password for your app stake keys. This password will be used to encrypt/decrypt your app-stake private keys stored in a Postgres database. It can be any combination of plaintext, symbols, letters, etc.

## 3. Configure Environment Variables
Fill out the `.env` variables for the gateway server. This can be done by injecting environment variables directly or using a `.env` file.

### Env Variables Description
| Variable Name                      | Description                            | Example Value                                      |
|------------------------------------|----------------------------------------|----------------------------------------------------|
| `POKT_RPC_FULL_HOST`               | Used for dispatching sessions          | `http://localhost:3000`                            |
| `HTTP_SERVER_PORT`                 | Gateway server port                    | `8080`                                             |
| `POKT_RPC_TIMEOUT`                 | Max response time for a Pokt node      | `5s`                                               |
| `ENVIRONMENT_STAGE`                | Log verbosity                          | `development`, `production`                        |
| `SESSION_CACHE_TTL`                | Duration for sessions to stay in cache | `75m`                                              |
| `POKT_APPLICATIONS_ENCRYPTION_KEY` | User-generated encryption key          | `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`                 |
| `DB_CONNECTION_URL`                | PostgreSQL Database connection URL     | `postgres://user:password@localhost:5432/postgres` |

See [.env.sample](..%2F.env.sample) for a sample.

## 4. Run Migration Script
Run the migration script to seed your PostgreSQL database.
```sh
./scripts/migration.sh -u
```

## 5. Insert App Stake Private Keys
Copy and paste the following SQL query to insert app stake private keys into the database:

```sql

INSERT INTO pokt_applications (encrypted_private_key)
VALUES (pgp_sym_encrypt('{private_key}', '{encryption_key}'));
```
_Note: Replace {private_key} and {encryption_key} in the SQL query with your actual private key and encryption key._

## 6. Compile and Run Gateway Server
Copy and paste the following code to compile and run the gateway server. Hit the endpoint http://localhost/relay/{chain_id} with a JSON-RPC payload.

```sh
go build cmd/gateway_server/main.go
```
_Remember to keep sensitive information secure and follow best practices for handling private keys and passwords._