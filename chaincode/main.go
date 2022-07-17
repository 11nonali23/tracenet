package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Campaign struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Recipient string `json:"recipient"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func (s *SmartContract) Test(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// Create a new campaign
func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, recipient string, startTime string, endTime string) error {
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since %s id already exists", id)
	}

	campaign := Campaign{
		ID:        id,
		Name:      name,
		Recipient: recipient,
		StartTime: startTime,
		EndTime:   endTime,
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, campaignJSON)

	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) ReadAllCampaigns(ctx contractapi.TransactionContextInterface) ([]*Campaign, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}

	var campaigns []*Campaign

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, &campaign)
	}

	resultsIterator.Close()

	return campaigns, nil
}

func (s *SmartContract) GetAvailableCampaings(ctx contractapi.TransactionContextInterface) ([]*Campaign, error) {
	currentTime := time.Now().Format(time.RFC3339)
	queryString := fmt.Sprintf(`{"selector":{"endTime":{"$gt": "%s"}}}`, currentTime)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var campaigns []*Campaign

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, &campaign)
	}

	resultsIterator.Close()

	return campaigns, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating campaign chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting campaign chaincode: %v", err)
	}
}
