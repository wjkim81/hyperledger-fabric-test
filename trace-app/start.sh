#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0
# This code is based on code written by the Hyperledger Fabric community. 
# Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/fabcar/startFabric.sh
#
# Exit on first error

set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

# launch network; 
cd ../basic-network

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d ca.example.com orderer.example.com peer0.org1.example.com couchdb


# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 tuna catches
docker-compose -f ./docker-compose.yml up -d cli

#docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n trace-app -c '{"function":"queryCow","Args":["002021864077"]}'
sleep 10
