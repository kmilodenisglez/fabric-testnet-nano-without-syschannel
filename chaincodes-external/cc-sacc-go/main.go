package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"log"
	"os"
)

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	// See chaincode.env.example
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}
	simpleAssetChaincode := new(SimpleAsset)

	//if err := shim.Start(simpleAssetChaincode); err != nil {
	//	fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	//}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      simpleAssetChaincode,
		TLSProps: shim.TLSProperties{
			Disabled:      true, // le decimos que la comunicacion va a ser con TLS
			Key:           nil,
			Cert:          nil,
			ClientCACerts: nil,
		},
	}

	fmt.Println("starting the chaincode on address: ", config.Address)

	if err := server.Start(); err != nil {
		log.Panicf("error starting asset-transfer-basic chaincode: %s", err)
	}

}
