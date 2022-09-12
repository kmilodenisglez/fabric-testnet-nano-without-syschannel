#!/usr/bin/env sh

# change by chaincode name
export CC_NAME=basic
# remember to modify CHAINCODE_ID by the one returned by the peer install
export CHAINCODE_ID=_chaincode_id_here_
export CHAINCODE_SERVER_ADDRESS=127.0.0.1:9999

# look for binaries in local samples /bin directory
export PATH="${PWD}"/../bin:"$PATH"
export FABRIC_CFG_PATH="${PWD}"/../config

export CHANNEL_NAME=mychannel
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=127.0.0.1:7051
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH="${PWD}"/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

# to query the chaincode (peer chaincode invoke and peer chaincode query)
export ORDERER_ADDRESS=127.0.0.1:6050
export ORDERER_TLS_CA=${CRYPTO_PATH}/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt