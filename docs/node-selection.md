# POKT Gateway Node Selection

## Node Selection System Architecture


## Session "Priming"

The gateway kit server under the hood will asynchronously "prime" sessions by calling a POKT full node for a session and cache the session for each app stake.
Caching is an optimization given that a session only changes periodically (4 blocks/60 mins) and prevents multiple round trips for each request we serve since session metadata is required per request. 

![gateway-server-session-cache.png](resources%2Fgateway-server-session-cache.png)

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

Some existing implementations of Checks can be found in:
1. [evm_data_integrity_check.go](..%2Finternal%2Fqos_node_registry%2Fchecks%2Fevm_data_integrity_check.go)
2. [evm_height_check.go](..%2Finternal%2Fqos_node_registry%2Fchecks%2Fevm_height_check.go)

### Adding custom QoS checks

Every custom check must conform to the `CheckJob` interface. The gateway server provides a base check:
```go
type Check struct {
	nextCheckTime time.Time
	nodeList      []*models.QosNode
	pocketRelayer pokt_v0.PocketRelayer
}
```
that developers should inherit. This base check provides a list of nodes to check, a time variable to determine when to run the check again, and a `PocketRelayer` that allows the developer to send requests to the nodes in the network.

Implementing custom QoS checks will be dependent on the chain or data source the developer is looking to support.  For example, the developer may want to send a request to a Solana node with a custom JSON-RPC method to see if the node is synced.
If the node is not synced, the developer can set a custom punishment through the various functions exposed in [qos_node.go](..%2Finternal%2Fqos_node_registry%2Fmodels%2Fqos_node.go), such as `SetTimeoutUntil` to punish the node.

Once the developer is finished implementing the CheckJob, they can enable the QoS check by initializing the newly created check into the `getEnabledJobs` function inside [qos_node_registry_service.go](..%2Finternal%2Fqos_node_registry%2Fqos_node_registry_service.go) and setting up a PR to be added into the official repository.

## Future Improvements

- Long term persistent results
    - Pros: More data to work with on determining if a node is healthy
    - Cons: Expensive, more complex logic, and can be punishing to new node operators
- Rolling up the results for long term storage & historical look back



