package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../crypto-config/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/Admin@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/Admin@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/tlsca/tlsca.org1.example.com-cert.pem"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var assetID = fmt.Sprintf("asset%d", time.Now().UnixNano()/1e6)

func main() {
	fmt.Println("=========== Fabric Gateway Go Example ===========")

	clientConn := newGrpcConnection()
	defer clientConn.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork("mychannel")
	contract := network.GetContract("basic")

	// InitLedger
	fmt.Println("\n--> Submit Transaction: InitLedger")
	_, err = contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(err)
	}
	fmt.Println("*** InitLedger committed")

	// CreateAsset
	fmt.Printf("\n--> Submit Transaction: CreateAsset %s\n", assetID)
	_, err = contract.SubmitTransaction("CreateAsset", assetID, "blue", "10", "Alice", "500")
	if err != nil {
		panic(err)
	}
	fmt.Println("*** CreateAsset committed")

	// ReadAsset
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset %s\n", assetID)
	result, err := contract.EvaluateTransaction("ReadAsset", assetID)
	if err != nil {
		panic(err)
	}
	fmt.Println("*** ReadAsset result:\n", formatJSON(result))
}

// gRPC connection to peer
func newGrpcConnection() *grpc.ClientConn {
	certPEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(err)
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(cert)
	creds := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	conn, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	return conn
}

// Identity from Admin
func newIdentity() *identity.X509Identity {
	certPEM, err := readFirstFile(certPath)
	if err != nil {
		panic(err)
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		panic(err)
	}
	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		panic(err)
	}
	return id
}

// Sign function from Admin key
func newSign() identity.Sign {
	keyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := identity.PrivateKeyFromPEM(keyPEM)
	if err != nil {
		panic(err)
	}
	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}
	return sign
}

// Read first file in folder
func readFirstFile(dir string) ([]byte, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in %s", dir)
	}
	return os.ReadFile(path.Join(dir, files[0].Name()))
}

// Pretty print JSON
func formatJSON(data []byte) string {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return string(data)
	}
	return out.String()
}
