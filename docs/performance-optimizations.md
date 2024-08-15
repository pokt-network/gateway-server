# Performance Optimizations

1. [FastHTTP](https://github.com/valyala/fasthttp) for both HTTP Client/Server
2. [FastJSON](https://github.com/pquerna/ffjson) for performant JSON Serialization and Deserialization
3. Lightweight Pocket Client

## Pocket Client Optimizations

We have implemented our own lightweight Pocket client to enhance speed and efficiency.

Leveraging the power of [FastHTTP](https://github.com/valyala/fasthttp) and [FastJSON](https://github.com/pquerna/ffjson), our custom client achieves remarkable performance gains.

Additionally, it has the capability to properly parse node runner's POKT errors properly given that the network runs diverse POKT clients (geomesh, leanpokt, their own custom client).

### Why It's More Efficient/Faster

1. **FastHTTP:** This library is designed for high-performance scenarios, providing a faster alternative to standard HTTP clients. Its concurrency-focused design allows our Pocket client to handle multiple requests concurrently, improving overall responsiveness.
2. **FastJSON:** The use of FastJSON ensures swift and efficient JSON serialization and deserialization. This directly contributes to reduced processing times, making our Pocket client an excellent choice for high-scale web traffic.
