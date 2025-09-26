package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the chaincode for asset management
type SmartContract struct {
	contractapi.Contract
}

// Asset defines the attributes of an asset
type Asset struct {
	ID             string `json:"ID"`             // Unique asset identifier
	Color          string `json:"color"`          // Asset color
	Size           int    `json:"size"`           // Asset size
	Owner          string `json:"owner"`          // Current asset owner
	AppraisedValue int    `json:"appraisedValue"` // Appraised asset value
}

// QueryResult is used to handle query results
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Asset `json:"Record"`
}

// InitLedger initializes the ledger with a set of predefined assets and log messages
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("\n[InitLedger] Initializing ledger with default assets...")

	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, asset := range assets {
		fmt.Printf("[InitLedger] Adding asset: %s\n", asset.ID)
		if err := s.putAsset(ctx, &asset); err != nil {
			return fmt.Errorf("[InitLedger] Failed to add asset %s: %v", asset.ID, err)
		}
	}

	fmt.Println("[InitLedger] Ledger successfully initialized.\n")
	return nil
}

// CreateAsset creates a new asset with educational log messages
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id, color string, size int, owner string, appraisedValue int) error {
	fmt.Printf("\n[CreateAsset] Creating new asset: %s\n", id)

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("[CreateAsset] Asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	if err := s.putAsset(ctx, &asset); err != nil {
		return fmt.Errorf("[CreateAsset] Failed to create asset %s: %v", id, err)
	}

	fmt.Printf("[CreateAsset] Asset %s successfully created.\n", id)
	return nil
}

// CreateAssetUsingStructParam creates an asset using a struct with log messages
func (s *SmartContract) CreateAssetUsingStructParam(ctx contractapi.TransactionContextInterface, asset *Asset) error {
	fmt.Printf("\n[CreateAssetUsingStructParam] Creating asset from struct: %s\n", asset.ID)

	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("[CreateAssetUsingStructParam] Asset %s already exists", asset.ID)
	}

	if err := s.putAsset(ctx, asset); err != nil {
		return fmt.Errorf("[CreateAssetUsingStructParam] Failed to create asset %s: %v", asset.ID, err)
	}

	fmt.Printf("[CreateAssetUsingStructParam] Asset %s successfully created.\n", asset.ID)
	return nil
}

// putAsset saves an asset into the world state with log messages
func (s *SmartContract) putAsset(ctx contractapi.TransactionContextInterface, asset *Asset) error {
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("[putAsset] Failed to serialize asset %s: %v", asset.ID, err)
	}

	if err := ctx.GetStub().PutState(asset.ID, assetJSON); err != nil {
		return fmt.Errorf("[putAsset] Failed to store asset %s: %v", asset.ID, err)
	}

	fmt.Printf("[putAsset] Asset stored/updated: %s\n", string(assetJSON))
	return nil
}

// ReadAsset retrieves an asset from the ledger with log messages
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	fmt.Printf("\n[ReadAsset] Reading asset: %s\n", id)

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("[ReadAsset] Failed to read asset %s: %v", id, err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("[ReadAsset] Asset %s does not exist", id)
	}

	var asset Asset
	if err := json.Unmarshal(assetJSON, &asset); err != nil {
		return nil, fmt.Errorf("[ReadAsset] Failed to deserialize asset %s: %v", id, err)
	}

	fmt.Printf("[ReadAsset] Asset successfully read: %s\n", string(assetJSON))
	return &asset, nil
}

// UpdateAsset updates an existing asset with log messages
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id, color string, size int, owner string, appraisedValue int) error {
	fmt.Printf("\n[UpdateAsset] Updating asset: %s\n", id)

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("[UpdateAsset] Asset %s does not exist", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	if err := s.putAsset(ctx, &asset); err != nil {
		return fmt.Errorf("[UpdateAsset] Failed to update asset %s: %v", id, err)
	}

	fmt.Printf("[UpdateAsset] Asset %s successfully updated.\n", id)
	return nil
}

// DeleteAsset deletes an asset with log messages
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	fmt.Printf("\n[DeleteAsset] Deleting asset: %s\n", id)

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("[DeleteAsset] Asset %s does not exist", id)
	}

	if err := ctx.GetStub().DelState(id); err != nil {
		return fmt.Errorf("[DeleteAsset] Failed to delete asset %s: %v", id, err)
	}

	fmt.Printf("[DeleteAsset] Asset %s successfully deleted.\n", id)
	return nil
}

// AssetExists checks if an asset exists with log messages
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("[AssetExists] Failed to read asset %s: %v", id, err)
	}
	return assetJSON != nil, nil
}

// TransferAsset transfers ownership of an asset with log messages
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	fmt.Printf("\n[TransferAsset] Transferring asset %s to new owner: %s\n", id, newOwner)

	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	if err := s.putAsset(ctx, asset); err != nil {
		return "", err
	}

	fmt.Printf("[TransferAsset] Old owner: %s, New owner: %s\n", oldOwner, newOwner)
	return oldOwner, nil
}

// GetAllAssets returns all assets from the ledger with log messages
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	fmt.Println("\n[GetAllAssets] Retrieving all assets from the ledger...")

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("[GetAllAssets] Failed to query ledger: %v", err)
	}
	defer resultsIterator.Close()

	var results []QueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[GetAllAssets] Error iterating results: %v", err)
		}

		var asset Asset
		if err := json.Unmarshal(queryResponse.Value, &asset); err != nil {
			return nil, fmt.Errorf("[GetAllAssets] Failed to deserialize asset: %v", err)
		}

		results = append(results, QueryResult{Key: queryResponse.Key, Record: &asset})
		fmt.Printf("[GetAllAssets] Found asset: %s\n", queryResponse.Key)
	}

	fmt.Println("[GetAllAssets] Query completed.\n")
	return results, nil
}
