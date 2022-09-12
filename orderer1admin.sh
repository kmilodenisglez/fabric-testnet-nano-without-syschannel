#!/usr/bin/env sh

if ! command -v peer version &> /dev/null
then
    # look for binaries in local samples /bin directory
    export PATH="${PWD}"/../bin:"$PATH"
fi

export FABRIC_CFG_PATH="${PWD}"

export FABRIC_LOGGING_SPEC=debug:cauthdsl,policies,msp,common.configtx,common.channelconfig=info
export ORDERER_ADMIN_LISTENADDRESS=127.0.0.1:9443
export OSN_TLS_CA_ROOT_CERT=["${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt]
export ADMIN_TLS_SIGN_CERT="${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
export ADMIN_TLS_PRIVATE_KEY="${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

# using osnadmin binary to handle the channel joining for the orderer

# inspect channels configuration in the orderer
osnadmin channel list -o 127.0.0.1:9443 --ca-file "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt --client-cert "${PWD}"/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/tls/client.crt --client-key "${PWD}"/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/tls/client.key

# join the orderer to mychannel
osnadmin channel join --channelID mychannel --config-block channel-artifacts/mychannel.block -o 127.0.0.1:9443 --ca-file "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt --client-cert "${PWD}"/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/tls/client.crt --client-key "${PWD}"/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/tls/client.key

# inspect mychannel channel configuration in the orderer
osnadmin channel list --channelID mychannel -o 127.0.0.1:9443 --ca-file "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt --client-cert "${PWD}"/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/tls/client.crt --client-key "${PWD}"/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/tls/client.key
