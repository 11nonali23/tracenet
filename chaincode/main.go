package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Campaign struct {
	ID        		string `json:"id"`
	Name      		string `json:"name"`
	Recipient 		string `json:"recipient"`
	StartTime 		string `json:"startTime"`
	EndTime			string `json:"endTime"`
	Owner			Owner  `json:"owner"`
	//TODO add viewers
}

type Owner struct {
	ID 					string `json:"ID"`
	KnowledgeGraph		string `json:"KnowledgeGraph"`
	privacyPreferences	string `json:"privacyPreferences"`
}

func (s *SmartContract) Test(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, recipient string, startTime string, endTime string) error {
	exists, err := s.CampaignExists(ctx, id)
	
	if err != nil {
		return err
	}
	
	if exists {
		return fmt.Errorf("Error while creating campaign: the campaign %s already exist", id)
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

func (s *SmartContract) DeleteCampaign(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.CampaignExists(ctx, id)
	
	if err != nil {
		return err
	}
	
	if !exists {
		return fmt.Errorf("Error while deleting campaign: the campaign %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) CampaignExists(ctx contractapi.TransactionContextInterface, campaignID string) (bool, error) {
	campaignBytes, err := ctx.GetStub().GetState(campaignID)

	if err != nil {
		return false, fmt.Errorf("Failed to read campaign %s from world state. %v", campaignID, err)
	}

	return campaignBytes != nil, nil
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

func (s *SmartContract) ShareKnowledgeGraph(ctx contractapi.TransactionContextInterface, campaignId string, ownerId string, knowledgeGraph string, privacyPreferences string) error {
	queryString := fmt.Sprintf(`{"selector":{"id":{"$eq": "%s"}}}`, campaignId)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return err
	}

	defer resultsIterator.Close()

	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return err
	}

	var campaign Campaign
	err = json.Unmarshal(queryResponse.Value, &campaign)
	if err != nil {
		return err
	}

	owner := Owner{
		ID:					ownerId,			
		KnowledgeGraph:		knowledgeGraph,	
		privacyPreferences:	privacyPreferences,
	}
	campaign.Owner = owner

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(campaignId, campaignJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) GetAvailableCampaings(ctx contractapi.TransactionContextInterface) (bool, error) {

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
