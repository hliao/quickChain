# Quick Chain

**Quick Chain** is a high-performance blockchain designed for **fast block production** and can be used as a **distributed database**.  
It provides low-latency consensus, flexible data storage, and scalable architecture suitable for both financial and non-financial applications.

---

## 🚀 Features

- ⚡ **Fast Block Production** — sub-second block times, optimized for low-latency use cases.
- 📦 **Distributed Database** — store structured or unstructured data on-chain with high availability.
- 🔒 **Secure Consensus** — Byzantine Fault Tolerant consensus ensuring data consistency across nodes.
- 🌐 **Scalable Network** — horizontally scalable node architecture.
- 🔧 **Developer Friendly** — gRPC/REST APIs, SDKs, and CLI tools for easy integration.

---

## 📂 Project Structure

```
quickChain/
├── app/                # ABCI application (CounterApp)
├── cmd/
│   └── quickchaincli/  # CLI main that runs the ABCI server
├── config/             # Node/app configuration (placeholder)
├── tests/              # Tests (placeholder)
├── types/              # Shared types (placeholder)
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/your-org/quickChain.git
cd quickChain
```

### 2. Run Guide (start CometBFT first)
```bash
# Install CometBFT
go install github.com/cometbft/cometbft/cmd/cometbft@v1.0.1

# Initialize a local node (once)
$(go env GOPATH)/bin/cometbft init

# Terminal A: start CometBFT (defaults to proxy_app tcp://127.0.0.1:26658)
$(go env GOPATH)/bin/cometbft start

# Terminal B: run the ABCI app (CounterApp)
go run ./cmd/quickchaincli

# Terminal C: send a few transactions via RPC
curl -s "http://localhost:26657/broadcast_tx_commit?tx=\"tx1\"" > /dev/null
curl -s "http://localhost:26657/broadcast_tx_commit?tx=\"tx2\"" > /dev/null
curl -s "http://localhost:26657/broadcast_tx_commit?tx=\"tx3\"" > /dev/null

# Query current count (base64 -> text)
curl -s "http://localhost:26657/abci_query?data=0x00" \
| jq -r '.result.response.value' | base64 --decode; echo

# (Optional) send 1000 tx sequentially (zsh/bash)
for i in {1..1000}; do curl -s "http://localhost:26657/broadcast_tx_commit?tx=\"tx$i\"" > /dev/null; done

# (Optional) send 1000 tx with concurrency 20
seq 1 1000 | xargs -P 20 -I{} sh -c 'curl -s "http://localhost:26657/broadcast_tx_commit?tx=\"tx{}\"" > /dev/null'
```

---

## 📖 Example Use Cases

- **Financial applications**: low-latency payment settlement
- **Gaming**: real-time asset tracking and state synchronization
- **IoT / Edge**: distributed data collection with fault tolerance
- **Enterprise**: decentralized log or record management

---

## 🔮 Roadmap

- [ ] Add WebAssembly (Wasm) smart contract support
- [ ] Cross-chain interoperability
- [ ] Sharding for horizontal scalability
- [ ] Advanced query language for on-chain data

---

## 🤝 Contributing

Contributions are welcome!  
Please open an issue or submit a pull request.

---

## 📜 License

[Apache 2.0](LICENSE)
