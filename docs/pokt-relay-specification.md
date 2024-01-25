# POKT Gateway Server Endpoints

## Overview
To truly understand why gateway operators and the Nodies Gateway Stack are important, let's review what steps application developers have to perform in order to send a single request to Pocket Network without gateways.

The chronological steps assuming the application is staked are:
* Generate an Application Authentication Token (AAT)
* Send a request to a Pocket full node for the latest nodes in a session
* Construct and sign a Relay Proof and submit it to one of the nodes in a session
* Receive a response from a node
* Determine if the response is legit or valid
* Proxy it back to your Web Application
---

## What are AATs?
The AAT is an auth token that allows application clients to access the network without the need to expose their private keys.

_Note: AATs are non-revocable and do not have a time expiration date. The only way to revoke a token is to unstake the entire application_

**AAT's Data Structure**
```go
type AAT struct {
   Version              string `protobuf:"bytes,1,opt,name=version,proto3" json:"version"`
   ApplicationPublicKey string `protobuf:"bytes,2,opt,name=applicationPublicKey,proto3" json:"app_pub_key"`
   ClientPublicKey      string `protobuf:"bytes,3,opt,name=clientPublicKey,proto3" json:"client_pub_key"`
   ApplicationSignature string `protobuf:"bytes,4,opt,name=applicationSignature,proto3" json:"signature"`
}
```
**JSON AAT Example**
```json
{
    "version": "0.0.1",
    "app_pub_key": 'eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743',
    "client_pub_key": "eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
    "signature": "5309f66a22ace63e8b4f94220151feabad11d4f3c22f50f6e395c72f1df96111da9bb25eceb11361d7e7074e7105d57dd2ec1d85cf962460608ef4bc5d35a80a",
}
```
---
## Generating an AAT
1. The AAT specification can be found here, but in simple terms:
2. JSON Encode AAT with an empty string signature field:
3. SHA3_256 the JSON bytes (bytes of the stringified JSON Object)
4. Sign with ED25519 ECDSA
5. HexEncode the result bytes into a string
6. Replace the empty AAT.signature field with the hex-encoded result.

An example code implementation of this can be found in both Javascript and Golang
1. [Javascript Implementation](https://github.com/pokt-network/pocket-aat-js/tree/staging)
2. [Golang Implementation](https://github.com/pokt-network/pocket-core/blob/staging/x/pocketcore/keeper/aat.go#L14)
---
## Retrieve the latest session for your application
An application is supported by 24 randomly chosen node runners for a short duration called a `session`, which currently lasts 4 blocks.

Once that session is over, the application gets a new set of node runners. To find out who the latest node runners are, the application connects to a Pocket node and sends it a request. This process is called dispatching.

Request (POST): `{pocket_host}/v1/client/dispatch` with the following payload:

```json
{
   "app_public_key":"514810e9139c5571905c642564b18cfb67899af2da05e638031075033da091a5",
   "chain":"0074"
}
```

Response:
```json
{
    "block_height": 108183,
    "session": {
        "header": {
            "app_public_key": "514810e9139c5571905c642564b18cfb67899af2da05e638031075033da091a5",
            "chain": "0074",
            "session_height": 108181
        },
        "key": "EKxfv3DhF8u7gn1dhZxjFPQFhE+FTGhjUtLCsnq6V4g=",
        "nodes": [
            {
                "address": "cd019c3b62cfb8cb9fd9863634fd42f2caef8984",
                "chains": [
                    "0021",
                    "0027",
                    "0052",
                    "0003",
                    "0006",
                    "0053",
                    "0004",
                    "0065",
                    "0070",
                    "0054",
                    "0048",
                    "0058",
                    "000F",
                    "0074",
                    "0009"
                ],
                "jailed": false,
                "output_address": "344e7bd9fc60a7f91f91b44219a5e7ef99af9810",
                "public_key": "fd332ff15904c5b6d68f42aa5f05c66e9f2d9ba8267014e80112a1a39105f4e8",
                "service_url": "https://6286.n.poktstaking.com:443",
                "status": 2,
                "tokens": "60010000000",
                "unstaking_time": "0001-01-01T00:00:00Z"
            },
```

_Note: Given that retrieving a session requires a full node, this means staked applications will need to source full nodes or run one themselves! As well, this acts as a failure point for app/gateway operators if they only rely on one full node for dispatching. Without a session, application developers cannot send a request. Thankfully thanks to pruning efforts and more full nodes entering the networks, this should become a lower risk._

---
## What are relay proofs?
At its core, a relay proof is like a digital receipt proving that an application sent a request to a node runner. Here's how it works:
* **Generation and Validation:** When an application makes a request, it creates and signs a 'relay proof'. This is like a digital signature, ensuring the request is genuine and hasn't been tampered with.
* **Verification by Servicer:** These servicers check the relay proof to make sure it's from a legitimate application in the network.
* **Storing the Proof:** Once verified, node runners store this proof in a data structure called a 'Merle Sum Index (MSI) Tree'.
* **Processing the Request:** The node runner then processes the request and sends the information back to the application.
* **Claim and Proof Lifecycle:** In the process of getting paid for their work, node runners go through a two-step 'claim and proof' cycle. First, they submit the 'root' of the MSI Tree as part of a claim transaction, indicating they have served several requests. Then, to provide evidence of their work, they must submit a randomly selected index along with Merkle proof from the branch to root, ensuring fairness and verification. This allows the network to trust that the node runners aren't just selecting an easy-to-prove transaction but deterministically chosen by the network in a secure way. Ultimately, this allows for a compute and space-efficient blockchain as validators of the network do not have to store nor verify every single request served.

## Generating a relay to the network
Now that we have access to a set of node runners (and assuming all node runners are actually operational), we still need to send the JSON-RPC request in a data structure that node runners will understand and accept. Unfortunately, node runners will not accept a simple HTTP JSON-RPC request as you would expect with other node providers, so we must construct a relay along with a relay proof.
POST request to `{pocket_host}/v1/client/relay` with the following payload
```json
{
   "payload":"relay_payload",
   "meta":"relay_meta",
   "proof":"relay_proof"
}
```

### **Relay Meta Data structure:**
```go
type Payload struct {
	Data    string            `json:"data"`              // the actual data string for the external chain
	Method  string            `json:"method"`            // the http CRUD method
	Path    string            `json:"path"`              // the REST Path
	Headers map[string]string `json:"headers,omitempty"` // http headers
}
```
```json
{
    "data": "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"0x050ea4ab4183E41129B7D72A492DaBf52B27EdB5\",\"latest\"],\"id\":67}",
    "method": "POST",
    "path": "",
    "headers": null
}
```

### **Relay Metadata structure:**
Relay metadata is simply when the request is made based on POKT's block time. This can be simplified to the session block height. However, if possible, the application developer should return the block height from /v1/query/height
```go
type RelayMeta struct {
   BlockHeight int64 `json:"block_height"` // the block height when the request is made
}
```
```json
{
    "block_height": 100
}
```

### **Relay Proof Data structure**
The generation of the relay proof is the most complex and crucial piece to ensuring the application's request is recognized and accepted by the node runners.
```json
type RelayProof struct {
   RequestHash        string `protobuf:"bytes,1,opt,name=requestHash,proto3" json:"request_hash"`
   Entropy            int64  `protobuf:"varint,2,opt,name=entropy,proto3" json:"entropy"`
   SessionBlockHeight int64  `protobuf:"varint,3,opt,name=sessionBlockHeight,proto3" json:"session_block_height"`
   ServicerPubKey     string `protobuf:"bytes,4,opt,name=servicerPubKey,proto3" json:"servicer_pub_key"`
   Blockchain         string `protobuf:"bytes,5,opt,name=blockchain,proto3" json:"blockchain"`
   Token              AAT    `protobuf:"bytes,6,opt,name=token,proto3" json:"aat"`
   Signature          string `protobuf:"bytes,7,opt,name=signature,proto3" json:"signature"`
}
```
1. Entropy Generate a random integer from [0, INT64]
2. `SessionBlockHeight`, `Blockchain`, `ServicerPubKey`, `Token` are retrievable from the above steps.
3. `RequestHash` is generated by with the following psuedo code
```go
// requestHash
{
    "payload": relay_payload,
    "meta": relay_meta
}
requestHashBytes := json.Marshal(requestHash) (bytes of the stringified JSON Object)
SHA3_256(requestHashBytes)
HexEncode to string
```
4. Once the relay object is generated, construct an ordered version of the relay object to hash and sign with the application private key for a signature using the following order:
   ng the following order:
```json
// relayProof
{
    "entropy": "1234567890123456",
    "session_block_height":  108181,
    "servicer_pub_key":  "a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890",
    "blockchain": "0074",
    "signature": "",
    "token": "SHA3_256-aat-without-signature",
    "request_hash": "sha-256-request-hash"
}
// json encode relay proof
relayProofJsonBytes := json.Marshal(relayProof)
hashedRelayProof = SHA3_256(rrelayProofJsonBytes)
SIGN(hashedRelayProof, appPrivateKey) -> c1d2e3f4a5b67890c1d2e3f4a5b67890c1d2e3f4a5b67890c1d2e3f4a5b67890
```
5. Use the generated signature to fill out the missing Signature field

**_NOTE: Ordering of the JSON object matters because the values are hashed. If the ordering changes, so will the hash._**

---
## Final Steps
With all the fields now generated, the valid relay proof can be constructed as below:
```json
// relay proof with signature
{
    "request_hash": "b1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "entropy": 1234567890123456,
    "session_block_height": 108181,
    "servicer_pub_key": "a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890",
    "blockchain": "0074",
    "aat": {
       ...
       // includes signature
    },
    "signature": "c1d2e3f4a5b67890c1d2e3f4a5b67890c1d2e3f4a5b67890c1d2e3f4a5b67890"
}

// Send to /v1/client/relay
{
   "payload":"relay_payload",
   "meta":"relay_meta",
   "proof":"relay_proof_with_signature"
}
```
If all goes well, the application should receive a response from the node runner!

----
After delving into the complexities of selecting a reliable source of dispatchers to retrieve a session, considering the network does not offer QoS assurances, and grasping the intricacies of sending requests to node runners, it becomes evident how crucial it is for software to abstract away the protocol and foster true developer adoption. This highlights the importance of the Gateway Operators and ultimately the Gateway server vision.