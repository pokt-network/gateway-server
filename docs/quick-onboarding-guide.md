# Quick Onboarding Guide

## 1. Preparing the Application Stakes
Application stakes provides Gateway Operators access to the POKT Network for sending traffic.
<details>
<summary>Testnet Instructions</summary>
<ol>
    <li>Generate 5 accounts (wallets) through the <a href="https://wallet.testnet.pokt.network">testnet wallet URL</a></li>
    <li>Distribute POKT to all the wallets generated through the <a href="https://faucet.pokt.network/">testnet faucet</a></li>
    <li>Stake each account into the network as an application stake with the chain id `0007` (a test chain that represents ETH Network). 
    <li>You can use our <a href="https://pokt-testnet-rpc.nodies.org">application stake script</a> to simplify the process if you don't have access to the staking CLI or the commands.</li>
</ol>
<hr>
Staking application stakes too complicated for you? 

No worries, we prestaked some shared applications stakes into POKT Testnet to help you get onboarded quicker. Please do not submit stake transactions to avoid disruption for other gateway operator testers as the applications are already staked on your behalf in the correct chain

All application  are staked into chain 0007 with 10M POKT.

<ul>
<li>1d06f04dcf5199a7f93f625d4fa507c2e0aca2f94fa3ebc2022c5e589406a9133d7ec4fef2ef676b340ce1df6ec5d0264ce1f40fae7fe9e07c415fa06fc1ffd6</li>
<li>2d0f9aab4396662db2a27d3388a1602e8081a49cb159471fdf4ef8aad4f9d120a1183ac69c10bf7f5df942b687b50a206fb1c54c66687c04c7710daed5f1e7a3</li>
<li>1e33f2948223e6655d4e10f462ad48203e18e81865098f4c15153ba4027f2fa4822fbcb6a0f485b9c61d1e84e976cb75214edc3e388b733e3ca4d5b80671cb4f</li>
<li>0bcdf221fb73f54a4acf4e61008a80c62ad155500846d99fd9cd190b46a9cf22157e1212fad906ac98bbf5a6b6ae50910ebd83e3fe789d3e4bd7f711abcd4ed1</li>
<li>20bf258e9e9632a9c627bfd328be87e0ecd6f14eeb7c7dc2382048c3063d3c08ec25b1aad594814f2a046cd2e89579992ecbba0951fec2d0f4b6ef1ba16fa8b9</li>
</ul>

</details>

<br> 
<details>
<summary>Mainnet Instructions</summary>
Application stakes in Morse are permissioned, therefore you must receive application stakes through the Pocket Network Foundation. If you are an authorized gateway operator, the Foundation will assist you in receiving the application stakes private keys.
</details>


## 2. Create Encryption Password
Create an encryption password for your app stake keys. This password will be used to encrypt/decrypt your app-stake private keys stored in a Postgres database. It can be any combination of plaintext, symbols, letters, etc.

## 3. Configure Environment Variables
Fill out the `.env` variables for the gateway server. This can be done by injecting environment variables directly or using a `.env` file.

### Env Variables Description
| Variable Name                      | Description                            | Example Value                                                                                     |
|------------------------------------|----------------------------------------|---------------------------------------------------------------------------------------------------|
| `POKT_RPC_FULL_HOST`               | Used for dispatching sessions          | `https://pokt-testnet-rpc.nodies.org` (a complimentary testnet dispatcher URL provided by Nodies) |
| `HTTP_SERVER_PORT`                 | Gateway server port                    | `8080`                                                                                            |
| `POKT_RPC_TIMEOUT`                 | Max response time for a Pokt node      | `5s`                                                                                              |
| `ENVIRONMENT_STAGE`                | Log verbosity                          | `development`, `production`                                                                       |
| `SESSION_CACHE_TTL`                | Duration for sessions to stay in cache | `75m`                                                                                             |
| `POKT_APPLICATIONS_ENCRYPTION_KEY` | User-generated encryption key          | `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`                                                                |
| `DB_CONNECTION_URL`                | PostgreSQL Database connection URL     | `postgres://user:password@localhost:5432/postgres`                                                |

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

If using testnet, send a request to chain `0007`, http://localhost/relay/0007. This testnet chain id represents Ethereum Node Operators and Nodies is currently supporting the chain id for reliable testing.

```sh
go build cmd/gateway_server/main.go
```
_Remember to keep sensitive information secure and follow best practices for handling private keys and passwords._