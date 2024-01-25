# POKT Gateway Node Selection

## Session Management

The gateway kit server under the hood will asynchronously call a POKT full node for a session and cache the session.
Caching is an optimization given that a session only changes periodically (4 blocks/60 mins).

![gateway-server-session-cache.png](resources%2Fgateway-server-session-cache.png)

## QoS Controls

The gateway kit server determines if a set of nodes are healthy based off a simple weight formula with the following
heuristics:

- Latency
- Success Responses
- Correctness in regard to other node operators.
- Liveliness (Synchronization)

```text
Weighted Score = Metrics.P90Latency * Weights["P90latency"] + Metrics.SuccessRate * Weights["successRate"] + Metrics.Correctness * Weights["correctness"]
```

The Gateway Server periodically will poll nodes in a session in order to score the nodes (fka as Nodies Cop), but will
also use real-time traffic to adjust the scores as well. Finally, the results will be stored in a short-term in-memory
cache which will allow for quick node selection.

## Default Thresholds and Weights

The gateway kit has default thresholds that are modifiable via the QoS table. In the event that there are no nodes that
are available for selection, it will default to a random node, then finally altruist.

**Threshholds:**

```json
"p90Latency":  300,
"successRate": 0.95,
"correctness": 0.95
"liveliness":  0.95,
```

**Weights:**

```json
"p90Latency":  0.05,
"successRate": 0.45,
"correctness": 0.50,
"liveliness":  -1
```

## Future Improvements

- Long term persistent results
    - Pros: More data to work with on determining if a node is healthy
    - Cons: Expensive, more complex logic, and can be punishing to new node operators
- Rolling up the results for long term storage & historical look back



