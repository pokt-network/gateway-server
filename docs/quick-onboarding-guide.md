# Mainnet Onboarding Guide
1. Retrieve your app stakes from the Foundation
2. Create an encryption password for your app stake keys (can be any plaintext, symbol, letters, etc).
   - This encryption password will be used to encrypt/decrypt your app-stake private keys stored in a Postgres database
3. Fill out the .env variables for the gateway server. 
   - This can be completed via injecting env variables directly or using a .env file
4. Run the migration script to seed your postgres db `./scripts/migration.sh -u`
5. Insert app stake private keys into the database with the below query:
```sql
INSERT INTO pokt_applications (encrypted_private_key)
VALUES (pgp_sym_encrypt('{private_key}', '{encryption_key}'));
```
6. Compile and run the gateway server and hit the endpoint http://localhost/relay/{chain_id} with a JSON-RPC payload.

## Env Variables Description
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