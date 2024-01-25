# POKT Gateway Server Overview

The Gateway server's goal is to reduce the complexity associated with directly interfacing with the protocol with an library that any developer can contribute to.  The Gateway server kickstarts off as a light-weight process that enables developers of all kinds to be able to interact with the protocol and engage with 50+ blockchains without the need to store terabytes of data, require heavy computational power, or understand the POKT protocol specifications with a simple docker-compose file. The rhetorical question that we pose to future actors who want to maintain a blockchain node is: Why spin up an Ethereum node and maintain it yourself whenever you can just leverage POKT natively using the Gateway server? After all, using POKT would require a fraction of the required resources and technical staffing.

## Features
- Simple docker-compose file with minimal dependencies to spin up
- a single tenancy HTTP endpoint for each blockchain that POKT supports, abstracting POKT's Relay Specification. This endpoint's throughput should scale based on the number of app stakes the developer provides.
- QoS checks to allow for optimized latency and success rates
- Provides Prometheus metrics for success, error rates, and latency for sending a relay to the network
- Custom Pocket client and web server that allows for efficient computational resources and memory footprint
- FastHTTP for optimized webserver and client
- FastJSON for efficient JSON Deserialization
- Custom Integration leveraging the two for efficient resource management
- Functionality improvement such as allowing for proper decoding of POKT Node error messages such as max evidence sealed errors.

## What's not included in the Gateway Server
- Authentication
- Rate Limiting & Multi-tenancy endpoints
- SaaS-based UX
- Reverse Proxy / Load Balancing Mechanisms
- Any other Opinionated SaaS-like design decision.

The decision to exclude certain features such as Authentication, Rate Limiting and multi-tenancy endpoints, SaaS-based UX, and Reverse Proxy/Load Balancing Mechanisms is rooted in the project's philosophy. These aspects are often regarded as opinionated web2 functionalities, and there are already numerous resources available on how to build SaaS products with various authentication mechanisms, rate-limiting strategies, and user experience design patterns. 

The Gateway server aims to simplify the POKT protocol, not reinventing the wheel. Each Gateway, being a distinct entity with its unique requirements and team dynamics, is better suited to decide on these aspects independently. For instance, the choice of authentication mechanisms can vary widely between teams, ranging from widely-used services like Auth0 and Amazon Cognito to in-house authentication solutions tailored to the specific language and skill set of the development team.

By not including these opinionated web2 functionalities, the Gateway server acknowledges the diversity of preferences and needs among developers and businesses. This approach allows teams to integrate their preferred solutions seamlessly, fostering flexibility and ensuring that the Gateway server remains lightweight and adaptable to a wide range of use cases.

As the project evolves, we anticipate that individual Gateways will incorporate their implementations of these features based on their unique requirements and preferences. This decentralized approach empowers developers to make decisions that align with their specific use cases, promoting a more customized and efficient integration with the Gateway server.

## Future
We envision that the server will be used as a foundation for the entirety of the ecosystem to continue to build on top of such as:
- Building their frontends and extending their backend to include POKT by using the gateway server for their own SaaS business
- Create Demo SaaS gateways that use the gateway server as the underlying foundation.
- Using POKT as a hyper scaler whenever they need more computational power or access to more blockchains (sticking the process into their LB rotation)
- Using POKT as a backend as a failover whenever their centralized nodes go down (sticking the process into their LB rotation)

Over time, as more gateways enter the network, there will be re-occurring patterns on what is needed on the foundational level and developers can create RFPs to have them included. For example, while rate limiting and multi-tenancy endpoints feel too opinionated right now, there is a future where we can create a service that distributes these endpoints natively in the GW server.  The use cases are limitless and we expect that over time, community contributions into the gateway server will enable some of the aforementioned use cases natively. 