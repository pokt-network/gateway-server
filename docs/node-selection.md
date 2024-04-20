# POKT Gateway Node Selection

## Node Selection System Architecture
![gateway-server-node-selection-system.png](resources%2Fgateway-server-node-selection-system.png)
- `Session Registry` - responsible for "priming" sessions asynchronously, providing session metadata, and feeding the node to the `NodeSelectorService`
- `Pocket Relayer` - responsible for sending a relay to the network
- `NodeSelectorService` - responsible for running QoS checks and identifying healthy nodes by chain.
- `ChainConfigurationRegistryService` - responsible for providing custom chain configurations such as altruists and timeouts.

## QoS Controls

The gateway kit server determines if a set of nodes are healthy based off a simple weight formula with the following
heuristics:

- Latency
- Success Responses
- Correctness in regard to other node operators.
- Liveliness (Synchronization)

## Node Selector
After the sessions are primed, the nodes are fed to the `NodeSelectorService` which is responsible for:
1. Running various QoS checks (Height and Data Integrity Checks)
2. Exposing functions for the main process to select a healthy node `findNode(chainId) string`

### Checks Framework
The gateway server provides a simple interface called a `CheckJob`. This interface consists of three simple functions
```go
type CheckJob interface {
  Perform()
  Name() string
  ShouldRun() bool
  }
```
Under the hood, the NodeSelectorService is responsible for asynchronously executing all the initialized `CheckJobs`.

### Existing QoS checks:
- **Height Check:** The general flow would be:
   1) Query all node operators height, 
   2) compares heights with other node operators within a specific threshold
   3) filters out node operators that exceed the configurable block height tolerance.

- **Data Integrity Check:** The general flow would be: 
   1) Retrieve a unique block identifier (i.e block hash or total block tx count, etc) with a configurable block offset for randomness, 
   2) Query other node operators for the same block identifier
   3) filter out other node operators that return a different identifier.

Some existing implementations of Checks can be found in:
1. [evm_data_integrity_check.go](..%2Finternal%2Fnode_selector_service%2Fchecks%2Fevm_data_integrity_check%2Fevm_data_integrity_check.go)
2. [evm_height_check.go](..%2Finternal%2Fnode_selector_service%2Fchecks%2Fevm_height_check%2Fevm_height_check.go)
3. [pokt_height_check.go](..%2Finternal%2Fnode_selector_service%2Fchecks%2Fpokt_height_check%2Fpokt_height_check.go)
4. [pokt_data_integrity_check.go](..%2Finternal%2Fnode_selector_service%2Fchecks%2Fpokt_data_integrity_check%2Fpokt_data_integrity_check.go)
5. [solana_height_check.go](..%2Finternal%2Fnode_selector_service%2Fchecks%2Fsolana_height_check%2Fsolana_height_check.go)
6. [solana_data_integrity_check.go](..%2Finternal%2Fnode_selector_service%2Fchecks%2Fsolana_data_integrity_check%2Fsolana_data_integrity_check.go)
### Adding custom QoS checks

Every custom check must conform to the `CheckJob` interface. The gateway server provides a base check:
```go
type Check struct {
NodeList           []*qos_models.QosNode
PocketRelayer      pokt_v0.PocketRelayer
ChainConfiguration chain_configurations_registry.ChainConfigurationsService
}
```
that developers should inherit. This base check provides a list of nodes to check and a `PocketRelayer` that allows the developer to send requests to the nodes in the network, and `ChainConfiguration` service that allows for per-chain specific check configurations.

Checks are designed to be opinionated and there are numerous ways to implement whether a node is healthy or not by definition. Therefore, implementing custom QoS checks will be dependent on the chain or data source the developer is looking to support.  For example, the developer may want to send a request to a custom blockchain node with a custom JSON-RPC method to see if the node is synced by using the provided `PocketRelayer` to send a request to the node through Pocket network.
If the node is not synced, the developer can set a custom punishment through the various functions exposed in [qos_node.go](..%2Finternal%2Fnode_selector_service%2Fmodels%2Fqos_node.go), such as `SetTimeoutUntil` to punish the node.

Once the developer is finished implementing the CheckJob, they can enable the QoS check by initializing the newly created check into the `enabledChecks` variable inside [node_selector_service.go](..%2Finternal%2Fnode_selector_service%2Fnode_selector_service.go) and are encouraged to open up a PR for inclusion in the official repository.

## Future Improvements

- Long term persistent results
    - Pros: More data to work with on determining if a node is healthy
    - Cons: Expensive, more complex logic (due to geographic regions) and can be punishing to new node operators
- Rolling up the results for long term storage & historical look back



