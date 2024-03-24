mockery --dir=./pkg/pokt/pokt_v0 --name=PocketService --filename=pocket_service_mock.go  --output=./mocks/pocket_service --outpkg=pocket_service_mock --with-expecter
mockery --dir=./pkg/ttl_cache --name=TTLCacheService --filename=ttl_cache_service_mock.go  --output=./mocks/ttl_cache_service --outpkg=ttl_cache_service_mock --with-expecter
mockery --dir=./internal/pokt_apps_registry --name=AppsRegistryService --filename=pokt_apps_registry_mock.go  --output=./mocks/apps_registry --outpkg=app_registry_mock --with-expecter
mockery --dir=./internal/session_registry --name=SessionRegistryService --filename=session_registry_mock.go  --output=./mocks/session_registry --outpkg=session_registry_mock --with-expecter
mockery --dir=./internal/chain_configurations_registry --name=ChainConfigurationsService --filename=chain_configurations_registry_mock.go  --output=./mocks/chain_configurations_registry --outpkg=chain_configurations_registry_mock --with-expecter
mockery --dir=./internal/node_selector_service --name=NodeSelectorService --filename=node_selector_mock.go  --output=./mocks/node_selector --outpkg=node_selector_mock --with-expecter
mockery --dir=./internal/apps_registry --name=AppsRegistryService --filename=app_registry_mock.go  --output=./mocks/apps_registry --outpkg=apps_registry_mock --with-expecter
mockery --dir=./internal/global_config --name=GlobalConfigProvider --filename=config_provider.go  --output=./mocks/global_config --outpkg=global_config_mock --with-expecter
