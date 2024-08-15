include .env

########################
### Makefile Helpers ###
########################

.PHONY: prompt_user
# Internal helper target - prompt the user before continuing
prompt_user:
	@echo "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help: ## Prints all the targets in all the Makefiles
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-60s\033[0m %s\n", $$1, $$2}'

########################
### Database Helpers ###
########################

.PHONY: db_migrate
db_migrate: ## Run database migrations
	@echo "Running database migrations..."
	./scripts/db_migrate.sh -u

.PHONY: db_migrate
db_migrate: ## Run database migrations
	@echo "Running database migrations..."
	./scripts/db_migrate.sh -u\


PG_CMD := INSERT INTO pokt_applications (encrypted_private_key) VALUES (pgp_sym_encrypt('$(POKT_APPLICATION_PRIVATE_KEY)', '$(POKT_APPLICATIONS_ENCRYPTION_KEY)'));
db_insert_app_private_key: ## Insert application private key into database
	@echo "Running SQL command..."
	@echo "$(PG_CMD)"
	@psql "$(DB_CONNECTION_URL)" -c "$(PG_CMD)"
