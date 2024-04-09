## Optimization Dependencies
- [FastHTTP](https://github.com/valyala/fasthttp) for both HTTP Client/Server
- [FastJSON](https://github.com/pquerna/ffjson) for performant JSON Serialization and Deserialization
- Lightweight Pocket Client

We have implemented our own lightweight Pocket client to enhance speed and efficiency. Leveraging the power of [FastHTTP](https://github.com/valyala/fasthttp) and [FastJSON](https://github.com/pquerna/ffjson), our custom client achieves remarkable performance gains. Additionally, it has the capability to properly parse node runner's POKT errors properly given that the network runs diverse POKT clients (geomesh, leanpokt, their own custom client).