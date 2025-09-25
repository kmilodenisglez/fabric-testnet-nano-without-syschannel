# Fabric test network - Nano bash - + Without a __system channel__ +

Based on fabric-samples/test-network-nano-bash, but without using __system channel__.

Use the same structure as test-network-nano-bash, see [test-network-nano-bash](https://github.com/hyperledger/fabric-samples/tree/main/test-network-nano-bash) for more details

# Pre-requirements

## Download fabric binaries <a name="fabconnect_testnetwork_download_prerequisites"></a>

```bash
mkdir fabric-folder && cd fabric-folder
```
To get the installation script:

```bash
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
```

Run the script:
```bash
./install-fabric.sh --fabric-version 2.5.12 b
```
> **NOTE**: These arguments download the `Fabric binaries`.

If you want to know more about the install-fabric.sh script visit the following link [install script](https://hyperledger-fabric.readthedocs.io/en/latest/install.html)

Clone this repo:
```bash
git clone https://github.com/kmilodenisglez/fabric-testnet-nano-without-syschannel.git
```

## To run the chaincode as a service
> **The (ccaas) builder release:** Since Fabric version 2.4.6 the chaincode-as-a-service (ccaas) builder release is available in all release distributions.

- If you are using a version prior to 2.4.6, and you do not have them in `fabric-folder/bin` you can build them from the Fabric source with the command `make ccaasbuilder`, you will then find the builder in `fabric/release/darwin-amd64/bin` or equivalent for your system. Just move the whole hierarchy starting there to `fabric-folder/bin` with something like: `mv release/darwin-amd64/bin/ccaas_builder ../fabric-folder/bin`

- You need to edit the `fabric-folder/config/core.yaml` file to point to that builder. The path specified in the default config file is only valid within the peer container which you won't be using. Modify the `externalBuilders` field in the `core.yaml` file to add the local external builder so that the configuration looks something like the following:
```
externalBuilders:
       - name: ccaas
         path: ../builders/ccaas
         propagateEnvironment:
           - CHAINCODE_AS_A_SERVICE_BUILDER_CONFIG
```
The path must be absolute or relative to where the peer will run so that it can find the builder when installing the chaincode.

# Instructions for starting network

👀 Note, by default you can start with a single ordering service node (OSN) and a single Org1 peer node and single Org1 peer admin terminal if you would like to keep things even more minimal (a single peer from Org1 can be utilized since the endorsement policy is set as any single organization).

> For several orderer nodes: Replace the `configtx.yaml` with `configtx-multiple-orderingnodes.yaml`.

Open terminal windows for 1 OSN, 1 peer node, and 1 peer admin as seen in the following terminal setup. The peer and peer admin belong to Org1.

![Terminal setup](terminal_setup.png)

The following instructions will have you run simple bash scripts that set environment variable overrides for a component and then runs the component.
The scripts contain only simple single-line commands so that they are easy to read and understand.
If you have trouble running bash scripts in your environment, you can just as easily copy and paste the individual commands from the script files instead of running the script files.

- cd to the `fabric-testnet-nano-without-syschannel` directory in each terminal window
- In the orderer terminal, run `./generate_artifacts.sh` to generate crypto material (calls cryptogen) and system and application channel genesis block and configuration transactions (calls configtxgen). The artifacts will be created in the `crypto-config` and `channel-artifacts` directories.
- Run `./orderer1.sh`
- In the admin terminal, run `./orderer1admin.sh`
- In the peer terminal, run `./peer1.sh`
- Note that each orderer and peer write their data (including their ledgers) to their own subdirectory under the `data` directory
- In the admin terminal, run `./peer1admin.sh`

The `peer1admin.sh` script sets the peer1 admin environment variables, creates the application channel `mychannel`, updates the channel configuration for the org1 gossip anchor peer, and joins peer1 to `mychannel`.
The remaining peer admin scripts join their respective peers to `mychannel`.

# Instructions for deploying and running the basic asset transfer sample go chaincode as a service

To deploy and invoke the chaincode, utilize the peer1 admin terminal that you have created in the prior steps.

## 1. Running the chaincode as a service

In another terminal, export the environment variables in the terminal:
```bash
source ./setenv.sh
```
> Note the syntax of running the scripts. The setenv.sh scripts run with the `source` command in order to source the script files in the respective shells. This is important so that the exported environment variables can be utilized by any subsequent user commands. In order to use commands like the `peer`.

Navigate to chaincode ex: `chaincodes-external/cc-assettransfer-go`:

```bash
cd chaincodes-external/cc-assettransfer-go
```

Package and install the external chaincode on peer1 with the following simple commands:

Package:
```bash
tar cfz code.tar.gz connection.json
tar cfz external-chaincode.tgz metadata.json code.tar.gz
```

Install the `cc-assettransfer-go` chaincode
```bash
peer lifecycle chaincode install external-chaincode.tgz
```

Run the following command to query all chaincode ID that you just installed:
```bash
peer lifecycle chaincode queryinstalled
```

The command will return output similar to the following:
```bash
Installed chaincodes on peer:
Package ID: basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80, Label: basic_1.0
```

Copy the returned chaincode package ID into an environment variable for use in subsequent commands (your ID may be different):

```bash
# in linux, wsl and darwin use export, ex:
export CHAINCODE_ID=basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80

# in windows terminal use set, ex:
set CHAINCODE_ID=basic_1.0:f3e2ca5115bba71aa2fd16e35722b420cb29c42594f0fdd6814daedbc2130b80
```

## Set chaincode name 
```bash
export CC_NAME=basic
```

## Approve and commit the chaincode

Using the peer1 admin, approve and commit the chaincode (only a single approve is required based on the lifecycle endorsement policy of any organization).

Approve chaincode:
```bash
peer lifecycle chaincode approveformyorg --version 1 --sequence 1 -o $ORDERER_ADDRESS --channelID $CHANNEL_NAME --name $CC_NAME --package-id $CHAINCODE_ID --tls --cafile $ORDERER_TLS_CA
```
Commit chaincode:
```bash
peer lifecycle chaincode commit --version 1 --sequence 1 -o $ORDERER_ADDRESS --channelID $CHANNEL_NAME --name $CC_NAME --tls --cafile $ORDERER_TLS_CA
```

Set the chaincode server address:
```bash
export CHAINCODE_SERVER_ADDRESS=127.0.0.1:9999

# windows
set CHAINCODE_SERVER_ADDRESS=127.0.0.1:9999
```

Build the chaincode:

```bash
# linux or darwin
go build -o ccass_binary

# windows
go build -o ccass_binary.exe
```

And start the chaincode service:

```bash
# linux
./ccass_binary

# windows
ccass_binary.exe
```

## Interact with the chaincode

In another terminal invoke the chaincode to create an asset (only a single endorser is required based on the default endorsement policy of any organization).
Then query the asset, update it, and query again to see the resulting asset changes on the ledger. Note that you need to wait a bit for invoke transactions to complete.

Export the environment variables in the terminal:
```bash
source ./setenv.sh
```

Populate the ledger with fake data:
### Init the ledger
```bash
peer chaincode invoke -c '{"Args":["InitLedger"]}' -o $ORDERER_ADDRESS -C $CHANNEL_NAME -n $CC_NAME --tls --cafile $ORDERER_TLS_CA
```

Insert an asset in ledger:
### Create an asset
```bash
peer chaincode invoke -c '{"Args":["CreateAsset","1","blue","35","tom","1000"]}' -o $ORDERER_ADDRESS -C $CHANNEL_NAME -n $CC_NAME --tls --cafile $ORDERER_TLS_CA
```

Query the asset created in the previous step:
### Read an asset
```bash
peer chaincode query -c '{"Args":["ReadAsset","1"]}' -C $CHANNEL_NAME -n $CC_NAME
```

Update the asset:
### Update an asset
```bash
peer chaincode invoke -c '{"Args":["UpdateAsset","1","blue","35","jerry","1000"]}' -o $ORDERER_ADDRESS -C $CHANNEL_NAME -n $CC_NAME --tls --cafile $ORDERER_TLS_CA
```

Query the asset modified in the previous step:
### Read an asset
```bash
peer chaincode query -c '{"Args":["ReadAsset","1"]}' -C $CHANNEL_NAME -n $CC_NAME
```

Congratulations, you have deployed a minimal Fabric network! Inspect the scripts if you would like to see the minimal set of commands that were required to deploy the network.

Utilize `Ctrl-C` in the orderer and peer terminal windows to kill the orderer and peer processes. You can run the scripts again to restart the components with their existing data, or run `./generate_artifacts` again to clean up the existing artifacts and data if you would like to restart with a clean environment.
