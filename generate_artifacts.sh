#!/usr/bin/env sh
set -eu

# remove existing artifacts, or proceed on if the directories don't exist
rm -r "${PWD}"/channel-artifacts || true
rm -r "${PWD}"/crypto-config || true
rm -r "${PWD}"/data || true

# look for binaries in local dev environment /build/bin directory and then in local samples /bin directory
export PATH="${PWD}"/../../fabric/build/bin:"${PWD}"/../bin:"$PATH"

echo "Generating MSP certificates using cryptogen tool"
cryptogen generate --config="${PWD}"/crypto-config.yaml

# set FABRIC_CFG_PATH to configtx.yaml directory that contains the profiles
export FABRIC_CFG_PATH="${PWD}"

echo "Generating  genesis block for application channel"
configtxgen -profile SampleAppChannelEtcdRaft -outputBlock channel-artifacts/mychannel.block -channelID mychannel

echo "Generating anchor peer update transaction for Org1"
configtxgen -profile SampleAppChannelEtcdRaft -outputAnchorPeersUpdate channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

echo "Generating anchor peer update transaction for Org2"
configtxgen -profile SampleAppChannelEtcdRaft -outputAnchorPeersUpdate channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
