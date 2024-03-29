#!/usr/bin/env sh

if ! command -v peer version &> /dev/null
then
    # look for binaries in local samples /bin directory
    export PATH="${PWD}"/../bin:"$PATH"
fi

export FABRIC_CFG_PATH="${PWD}"/../config

export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=127.0.0.1:7051
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH="${PWD}"/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

# peer1 admin will be responsible for adding anchor peer
#peer channel update -o 127.0.0.1:6050 -c mychannel -f "${PWD}"/channel-artifacts/Org1MSPanchors.tx --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt

# join peer to channel
peer channel join -b "${PWD}"/channel-artifacts/mychannel.block
