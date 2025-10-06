 Hyperledger Fabric Assignment - Golang Implementation

## Project structure
```
fabric-assignment/
├── chaincode/asset/asset.go
├── api/
│   ├── server.go
│   ├── go.mod
│   └── Dockerfile
└── docs/
    └── Assignment_Documentation.pdf
```

## What is included
- Go chaincode `asset.go` implementing CreateAsset, ReadAsset, UpdateBalance, GetHistory.
- Go REST API `server.go` that connects to Fabric Gateway (placeholder config).
- Dockerfile to containerize the REST API.
- PDF documentation with step-by-step instructions.

## Important notes
- This repository contains code files only. It does NOT include Fabric binaries, docker images, or a running network.
- To run the full solution you must have:
  - Docker and Docker Compose
  - Hyperledger Fabric binaries (peer, orderer, configtxgen, etc.)
  - Fabric test-network (from `fabric-samples/test-network`)
  - Go (>=1.20)

## Quick start (high level)
1. Download and install Fabric samples and binaries (on a machine with internet):
   ```bash
   curl -sSL https://bit.ly/2ysbOFE | bash -s
   ```
2. Start the test network:
   ```bash
   cd fabric-samples/test-network
   ./network.sh up createChannel -c mychannel -ca
   ```
3. Deploy chaincode:
   ```bash
   # copy chaincode to fabric-samples/chaincode/asset and then:
   ./network.sh deployCC -ccn asset -ccp ../chaincode/asset -ccl go
   ```
4. Prepare gateway connection and wallet (see PDF docs).
5. Build and run REST API:
   ```bash
   cd /path/to/fabric-assignment/api
   go mod tidy
   go run server.go
   # or build docker image
   docker build -t asset-api .
   docker run -p 3000:3000 --env-file .env asset-api
   ```
   
  



