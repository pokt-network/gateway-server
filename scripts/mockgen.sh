mockery --dir=./pkg/pokt/pokt_v0 --name=PocketService --filename=pocket_service_mock.go  --output=./mocks --outpkg=mocks --with-expecter
mockery --dir=./pkg/ttl_cache --name=TTLCacheService --filename=ttl_cache_service_mock.go  --output=./mocks --outpkg=mocks --with-expecter
mockery --dir=./internal/pokt_apps_registry --name=AppsRegistryService --filename=pokt_apps_registry_mock.go  --output=./mocks --outpkg=mocks --with-expecter
mockery --dir=./internal/session_registry --name=SessionRegistryService --filename=session_registry_mock.go  --output=./mocks --outpkg=mocks --with-expecter
mockery --dir=./internal/altruist_registry --name=AltruistRegistryService --filename=altruist_registry_mock.go  --output=./mocks --outpkg=mocks --with-expecter