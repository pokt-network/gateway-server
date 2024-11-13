<div align="center">
<img src="https://github.com/pokt-network/gateway-server/blob/main/docs/resources/gateway-server-logo.jpg" width="500" alt="POKT Gateway Server">
</div>

> [!WARNING]
>
> 🚧🚧🚧 This repository is `Archived`! As of November 13, 2024 this repository is `Archived`. We encourage users to migrate and adopt the [Path API and Toolkit Harness (PATH)](https://github.com/buildwithgrove/path) 🚧🚧🚧

# What is POKT Gateway Server? <!-- omit in toc -->

_tl;dr Streamline access to POKT Network's decentralized supply network._

The POKT Gateway Server is a comprehensive solution designed to simplify the integration of applications with POKT Network. Its goal is to reduce the complexities associated with directly interfacing with the protocol, making it accessible to a wide range of users, including application developers, existing centralized RPC platforms, and future gateway operators.

Learn more about the vision and overall architecture [overview](./docs/overview.md).

- [Gateway Operator Quickstart Guide](#gateway-operator-quickstart-guide)
  - [Interested in learning more?](#interested-in-learning-more)
- [Docker Image Releases](#docker-image-releases)
- [Docker Compose](#docker-compose)
- [Minimum Hardware Requirements](#minimum-hardware-requirements)
- [Database Migrations](#database-migrations)
  - [Creating a DB Migration](#creating-a-db-migration)
  - [Applying a DB Migration](#applying-a-db-migration)
  - [DB Migration helpers](#db-migration-helpers)
    - [Applying Migrations](#applying-migrations)
    - [Migrations Rollbacks](#migrations-rollbacks)
- [Unit Testing](#unit-testing)
  - [Generating Mocks](#generating-mocks)
  - [Running Tests](#running-tests)
- [Generating DB Queries](#generating-db-queries)
- [Contributing Guidelines](#contributing-guidelines)
- [Project Structure](#project-structure)

## Gateway Operator Quickstart Guide

To onboard the gateway server without having to dig deep, you can follow the [Quick Onboarding Guide](docs/quick-onboarding-guide.md).

### Interested in learning more?

We have an abundance of information in the [docs](docs) section:

1. [Gateway Server Overview](docs/overview.md)
2. [Gateway Server API Endpoints](docs/api-endpoints.md)
3. [Gateway Server System Architecture](docs/system-architecture.md)
4. [Gateway Server Node Selection](docs/node-selection.md)
5. [POKT Primer](docs/pokt-primer.md)
6. [POKT's Relay Specification](docs/pokt-relay-specification.md)

## Docker Image Releases

Every release candidate is published to [gateway-server/pkgs/container/pocket-gateway-server](https://github.com/pokt-network/gateway-server/pkgs/container/pocket-gateway-server).

## Docker Compose

There is an all-inclusive docker-compose file available for development [docker-compose.yml](docker-compose.yml.sample)

## Minimum Hardware Requirements

To run a Gateway Server, we recommend the following minimum hardware requirements:

- 1GB of RAM
- 1GB of storage
- 4 vCPUs+

In production, we have observed memory usage increase to 4GB+. The memory footprint will be dependent on the number of app stakes/chains staked and total traffic throughput.

## Database Migrations

<!-- TODO_IMPROVE: Need more details on why we need a DB, how it's used, pre-requisites, etc... -->

### Creating a DB Migration

Migrations are like version control for your database, allowing your team to define and share the application's database schema definition.

Before running a migration make sure to install the go lang migration cli on your machine. See [golang-migrate/migrate/tree/master/cmd/migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) for reference.

```sh
./scripts/migration.sh -n {migration_name}
```

This command will generate a up and down migration in `db_migrations`

### Applying a DB Migration

DB Migrations are applied upon server start, but as well, it can be applied manually through:

```sh
./scripts/migration.sh {--down or --up} {number_of_times}
```

### DB Migration helpers

#### Applying Migrations

- To apply all migrations:

  ```sh
  ./scripts/migration.sh --up
  ```

- To apply a specific number of migrations:

  ```sh
  ./scripts/migration.sh --up 2
  ```

#### Migrations Rollbacks

Make sure to provide either the number of migrations to rollback or the `--all` flag to rollback all migrations.

- To roll back a specific number of migrations:

  ```sh
  ./scripts/migration.sh --down 2
  ```

- To roll back all migrations:

  ```sh
  ./scripts/migration.sh --down --all
  ```

## Unit Testing

### Generating Mocks

Install Mockery with

```bash
go install github.com/vektra/mockery/v2@v2.40.1
```

You can generate the mock files through:

```sh
./scripts/mockgen.sh
```

By running this command, it will generate the mock files in `./mocks` folder.

Reference for mocks can be found [here](https://vektra.github.io/mockery/latest).

### Running Tests

Run this command to run tests:

```sh
go test -v -count=1  ./...
```

## Generating DB Queries

Gateway server uses [PGGen](https://github.com/jschaf/pggen) to create autogenerated type-safe queries.
Queries are added inside [queries.sql](./internal/Fdb_query/queries.sql) and re-generated via `./scripts/querygen.sh`.

## Contributing Guidelines

1. Create a Github Issue on the feature/issue you're working on.
2. Fork the project
3. Create new branch with `git checkout -b "branch_name"` where branch name describes the feature.
   - All branches should be based off `main`
4. Write your code
5. Make sure your code lints with `go fmt ./...` (This will Lint and Prettify)
6. Commit code to your branch and issue a pull request and wait for at least one review.
   - Always ensure changes are rebased on top of main branch.

## Project Structure

A partial high-level view of the code structure (generated)

```bash
.
├── cmd # Contains the entry point of the binaries
│   └── gateway_server # HTTP Server for serving requests
├── internal # Shared internal folder for all binaries
├── pkg # Distributable dependencies
└── scripts # Contains scripts for development
```

_Generate via `tree -L 2`_

---
