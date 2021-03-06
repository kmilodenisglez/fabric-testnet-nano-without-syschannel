# Fabric test network - Nano bash - + Without a __system channel__ +

Based on fabric-samples/test-network-nano-bash, but without using __system channel__.

Use the same structure as test-network-nano-bash, see [test-network-nano-bash](https://github.com/hyperledger/fabric-samples/tree/main/test-network-nano-bash) for more details

# Prereqs

- Follow the Fabric documentation for the [Prereqs](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)
- Follow the Fabric documentation for [downloading the Fabric samples and binaries](https://hyperledger-fabric.readthedocs.io/en/latest/install.html). You can skip the docker image downloads by using `curl -sSL https://bit.ly/2ysbOFE | bash -s -- -d`

## To run the chaincode as a service
- You need to have the `ccaas_builder` binaries. If you do not have them in `fabric-samples/bin` you can build them from the Fabric source with the command `make ccaasbuilder`, you will then find the builder in `fabric/release/darwin-amd64/bin` or equivalent for your system. Just move the whole hierarchy starting there to `fabric-samples/bin` with something like: `mv release/darwin-amd64/bin/ccaas_builder ../fabric-samples/bin`
- You need to edit the `fabric-samples/config/core.yaml` file to point to that builder. The path specified in the default config file is only valid within the peer container which you won't be using. Modify the `externalBuilders` field in the `core.yaml` file to add the local external builder so that the configuration looks something like the following:
```
externalBuilders:
       - name: ccaas_builder
         path: ../bin/ccaas_builder
         propagateEnvironment:
           - CHAINCODE_AS_A_SERVICE_BUILDER_CONFIG
       - name: other_ccaas_builder
         path: /opt/hyperledger/ccaas_builder
         propagateEnvironment:
           - CHAINCODE_AS_A_SERVICE_BUILDER_CONFIG
```
The path must be absolute or relative to where the peer will run so that it can find the builder when installing the chaincode.

# Instructions for starting network

## with multiple ordering nodes

Antes debe remplazar el contenido del fichero `configtx.yaml` por `configtx-multiple-orderingnodes.yaml`. 

Open terminal windows for 3 ordering nodes, 4 peer nodes, and 4 peer admins as seen in the following terminal setup. The first two peers and peer admins belong to Org1, the latter two peer and peer admins belong to Org2.
Note, you can start with two ordering nodes and a single Org1 peer node and single Org1 peer admin terminal if you would like to keep things even more minimal (two ordering nodes are required to achieve consensus (2 of 3), while a single peer from Org1 can be utilized since the endorsement policy is set as any single organization).
![Terminal setup](terminal_setup.png)

The following instructions will have you run simple bash scripts that set environment variable overrides for a component and then runs the component.
The scripts contain only simple single-line commands so that they are easy to read and understand.
If you have trouble running bash scripts in your environment, you can just as easily copy and paste the individual commands from the script files instead of running the script files.

- cd to the `test-network-nano-bash` directory in each terminal window
- In the first orderer terminal, run `./generate_artifacts.sh` to generate crypto material (calls cryptogen) and system and application channel genesis block and configuration transactions (calls configtxgen). The artifacts will be created in the `crypto-config` and `channel-artifacts` directories.
- In the three orderer terminals, run `./orderer1.sh`, `./orderer2.sh`, `./orderer3.sh` respectively
- In the three orderer admin terminals, run `./orderer1admin.sh`, `./orderer1admin.sh`, `./orderer1admin.sh` respectively
- In the four peer terminals, run `./peer1.sh`, `./peer2.sh`, `./peer3.sh`, `./peer4.sh` respectively
- Note that each orderer and peer write their data (including their ledgers) to their own subdirectory under the `data` directory
- In the four peer admin terminals, run `source peer1admin.sh`, `source peer2admin.sh`, `source peer3admin.sh`, `source peer4admin.sh` respectively

Note the syntax of running the scripts. The peer admin scripts run with the `source` command in order to source the script files in the respective shells. This is important so that the exported environment variables can be utilized by any subsequent user commands.

The `peer1admin.sh` script sets the peer1 admin environment variables, creates the application channel `mychannel`, updates the channel configuration for the org1 gossip anchor peer, and joins peer1 to `mychannel`.
The remaining peer admin scripts join their respective peers to `mychannel`.

# Instructions for deploying and running the basic asset transfer sample go chaincode as a service

To deploy and invoke the chaincode, utilize the peer1 admin terminal that you have created in the prior steps. You have two possibilities:

## 1. Running the chaincode as a service

Package and install the external chaincode on peer1 with the following simple commands:

```
cd chaincode-external

tar cfz code.tar.gz connection.json
tar cfz external-chaincode.tgz metadata.json code.tar.gz

cd ..

source ./setenv.sh

peer lifecycle chaincode install chaincode-external/external-chaincode.tgz
```

Run the following command to query the package ID of the chaincode that you just installed:
```
peer lifecycle chaincode queryinstalled -o 127.0.0.1:6050 --tlsRootCertFiles ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

The command will return output similar to the following:
```
Installed chaincodes on peer:
Package ID: basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80, Label: basic_1.0
```

Copy the returned chaincode package ID into an environment variable for use in subsequent commands (your ID may be different):

```
export CHAINCODE_ID=basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80
```

In another terminal, navigate to chaincode ex: `fabric-samples/asset-transfer-basic/chaincode-external` and build the chaincode:

```
# linux or darwin
go build -o ccass_binary

# windows
go build -o ccass_binary.exe
```

Set the chaincode package ID again (this is a different terminal):

```
export CHAINCODE_ID=basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80

# in windows use set, ex:

set CHAINCODE_ID=basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80
```

Set the chaincode server address:

```
export CHAINCODE_SERVER_ADDRESS=127.0.0.1:9999

# windows
set CHAINCODE_SERVER_ADDRESS=127.0.0.1:9999
```

And start the chaincode service:

```
# linux
./ccass_binary

# windows
ccass_binary.exe
```

## Activate the chaincode

Using the peer1 admin, approve and commit the chaincode (only a single approver is required based on the lifecycle endorsement policy of any organization):

```
peer lifecycle chaincode approveformyorg -o 127.0.0.1:6050 --channelID mychannel --name basic --version 1 --package-id $CHAINCODE_ID --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt

peer lifecycle chaincode commit -o 127.0.0.1:6050 --channelID mychannel --name basic --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

## Interact with the chaincode

Invoke the chaincode to create an asset (only a single endorser is required based on the default endorsement policy of any organization).
Then query the asset, update it, and query again to see the resulting asset changes on the ledger. Note that you need to wait a bit for invoke transactions to complete.

```
peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n basic -c '{"Args":["CreateAsset","1","blue","35","tom","1000"]}' --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt

peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","1"]}'

peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n basic -c '{"Args":["UpdateAsset","1","blue","35","jerry","1000"]}' --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt

peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","1"]}'
```

Congratulations, you have deployed a minimal Fabric network! Inspect the scripts if you would like to see the minimal set of commands that were required to deploy the network.

Utilize `Ctrl-C` in the orderer and peer terminal windows to kill the orderer and peer processes. You can run the scripts again to restart the components with their existing data, or run `./generate_artifacts` again to clean up the existing artifacts and data if you would like to restart with a clean environment.
