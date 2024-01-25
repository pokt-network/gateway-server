# POKT Primer

##  POKT Network: A Quick Overview
1. Apps (App developers): Individuals or entities that stake into the Pocket Network and obtain access to external blockchain nodes in return.
2. Node Runners: Individuals or entities that stake into the network and provide access to external blockchain nodes, such as Ethereum & Polygon in return for $POKT.

### The Challenge: Interacting with the Protocol
For application developers, directly engaging with the POKT Network can be intimidating due to its inherent complexities. The technical barriers and protocol nuances can dissuade many from integrating and adopting the network.

**Challenges faced by App Developers using the protocol:**
1. Managing Throughput: The network supports approximately 250 app stakes and around 10B relays. With each app stake being limited to roughly 20M requests per hour, developers who surpass this need to stake multiple applications and balance the load among these stakes.
2. Determining Quality of Service (QoS): The network doesn't currently enforce QoS standards. Apps are assigned a set of pseudo-randomly selected node runners, rotated over specified intervals, also known as sessions. It falls on the application developer to implement strategies, such as filtering and predictive analysis, to select node runners that align with their criteria for reliability, availability, and data integrity.
3. Protocol Interaction: Unlike the straightforward procedure of sending requests to an HTTP JSON-RPC server, interacting with the POKT Network requires far more complexities given its blockchain nature. (i.e. signing a request for a relay proof) 

### The Solution: Gateway Operators
Gateway Operators act as a conduit between app developers and the POKT Network, streamlining the process by abstracting the network's complexities. Their operations on a high level can be seen as:
1. Managing Throughput: By staking and load-balancing app stakes, Gateway Operators ensure the required throughput for smooth network interactions.
2. Determining the Quality of Service: Gateway Operators filter malicious, out-of-sync, offline, or underperforming nodes.
3. Protocol Interaction: Gateway Operators offer a seamless HTTP JSON-RPC Server interface, making it simpler for developers to send and receive requests, akin to interactions with conventional servers. Under the hood, the web server will contain the necessary business logic in order to interact with the protocol.

### Conclusion
Engaging with the POKT Network's capabilities doesn't have to be an uphill task. Thanks to Gateway Operators, app developers can concentrate on their core competencies—developing remarkable applications using a familiar HTTP interface, like traditional RPC providers—all while reaping the benefits of a decentralized RPC platform.

---

###### Footnotes:

1. _As of 9/14/2023, the app stakes are permissioned and overseen by the Pocket Network Foundation for security considerations._
2. _The amount of POKT staked into an app doesn't carry significant implications as all gateway operators are charged a fee for every request sent through an app stake._
3. _Historically, Grove (formerly known as Pocket Network Inc.) has been the sole gateway operator. This will change by 2024 Q2 as more gateway operators join the network._
4. _Our research aims to invite more gateway operators to join the network in a sustainable fashion by documenting the protocol specifications and limitations and leveraging and providing open-source software and noncloud vendor-lock-in services._