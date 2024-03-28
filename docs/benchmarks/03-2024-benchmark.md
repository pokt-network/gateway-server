# March 2024 Benchmark (RC 0.1.0) (WIP)

## Benchmark Purpose
The purpose of this benchmark is to comprehensively assess the performance metrics, particularly CPU and Memory behaviors, incurred while serving requests through the gateway server. Specifically, this evaluation aims to gauge the efficiency of various operations involved with sending a relay such as JSON serialization, IO handling, cryptographic signing, and asynchronous background processes (QoS checks).

## Benchmark Environment
- **POKT Testnet**: The benchmark is conducted within the environment of the POKT Testnet to simulate real-world conditions accurately.
- **RPC Method: eth_chainId**: The benchmark uses a consistent time RPC method, eth_chainId, chosen deliberately to isolate the impact of the gateway server overhead. It's noteworthy that the computational overhead of sending a request to the POKT Network remains independent of the specific RPC Method employed, hence why `eth_chainId` is chosen.
- **Hardware**: The gateway server is deployed on a dedicated instance, denoted by [x], ensuring controlled conditions for performance evaluation.
- **Tooling**: Utilizes Vegeta, a versatile HTTP load testing too

## Results
[Results section will provide detailed insights into the observed metrics, including CPU utilization, memory consumption, throughput, and any other pertinent performance indicators.]

## Summary
[Summary section will encapsulate key findings and conclusions drawn from the benchmarking exercise, highlighting any noteworthy trends, optimizations, or areas for improvement identified during the evaluation process.]
