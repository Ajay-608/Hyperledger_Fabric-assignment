# Hyperledger Fabric Assignment - Golang Implementation

## Project structure
-chaincode

-asset -> asset.go

-api

   - server.go,go.mod,Dockerfile


## What is included
- Go chaincode `asset.go` implementing CreateAsset, ReadAsset, UpdateBalance, GetHistory.
- Go REST API `server.go` that connects to Fabric Gateway (placeholder config).
- Dockerfile to containerize the REST API.
  

## Important notes
- This repository contains code files only. It does NOT include Fabric binaries, docker images, or a running network.
- To run the full solution you must have:
  - Docker and Docker Compose
  - Hyperledger Fabric binaries (peer, orderer, configtxgen, etc.)
  - Fabric test-network (from `fabric-samples/test-network`)
  - Go (>=1.20)



