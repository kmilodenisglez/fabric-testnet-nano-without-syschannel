# Sample smart contracts

This folder contains example smart contracts. It is recommended that users start with the Asset transfer samples and
tutorials series for the most recent example smart contracts.

| **Smart Contract** | **Description** | **Languages** |
|--------------------|------------------------------|---------|
| [cc-sacc-go](cc-sacc-go)  | Simple asset chaincode that interacts with the ledger using the low-level APIs provided by the Fabric Chaincode Shim API. | Go |


## Interact with the chaincode cc-sacc-go

## GET
peer chaincode invoke -c '{"Args":["get", "key1"]}' -o $ORDERER_ADDRESS -C $CHANNEL_NAME -n $CC_NAME --tls --cafile $ORDERER_TLS_CA

## SET
peer chaincode invoke -c '{"Args":["set", "key1", "value1"]}' -o $ORDERER_ADDRESS -C $CHANNEL_NAME -n $CC_NAME --tls --cafile $ORDERER_TLS_CA