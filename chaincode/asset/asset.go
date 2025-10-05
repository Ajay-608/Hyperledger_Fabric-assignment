package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Asset represents an account/asset in the ledger
type Asset struct {
	DealerID    string  `json:"dealerId"`
	MSISDN      string  `json:"msisdn"`
	MPIN        string  `json:"mpin"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	TransAmount float64 `json:"transAmount"`
	TransType   string  `json:"transType"`
	Remarks     string  `json:"remarks"`
}

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// CreateAsset adds a new asset to the world state
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerId, msisdn, mpin, balanceStr, status, transAmountStr, transType, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("asset %s already exists", dealerId)
	}

	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid balance: %v", err)
	}
	transAmount, err := strconv.ParseFloat(transAmountStr, 64)
	if err != nil {
		return fmt.Errorf("invalid transAmount: %v", err)
	}

	asset := Asset{
		DealerID:    dealerId,
		MSISDN:      msisdn,
		MPIN:        mpin,
		Balance:     balance,
		Status:      status,
		TransAmount: transAmount,
		TransType:   transType,
		Remarks:     remarks,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerId, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerId string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("asset %s does not exist", dealerId)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// UpdateBalance updates only the balance field of an asset
func (s *SmartContract) UpdateBalance(ctx contractapi.TransactionContextInterface, dealerId, newBalanceStr string) error {
	asset, err := s.ReadAsset(ctx, dealerId)
	if err != nil {
		return err
	}
	newBalance, err := strconv.ParseFloat(newBalanceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid newBalance: %v", err)
	}
	asset.Balance = newBalance
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(dealerId, assetJSON)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dealerId string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerId)
	if err != nil {
		return false, err
	}
	return assetJSON != nil, nil
}

// GetHistory returns the history of an asset
func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, dealerId string) ([]map[string]interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerId)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var history []map[string]interface{}
	for resultsIterator.HasNext() {
		mod, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var item map[string]interface{}
		if mod.IsDelete {
			item = map[string]interface{}{
				"txId":      mod.TxId,
				"isDelete":  mod.IsDelete,
				"timestamp": mod.Timestamp,
				"value":     nil,
			}
		} else {
			_ = json.Unmarshal(mod.Value, &item)
			item["txId"] = mod.TxId
			item["isDelete"] = mod.IsDelete
			item["timestamp"] = mod.Timestamp
		}
		history = append(history, item)
	}
	return history, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode: %v", err))
	}
	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}
