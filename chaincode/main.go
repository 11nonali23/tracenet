package main

import (
	"encoding/json"
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
	StartTime 		string `json:"startTime"`
	EndTime			string `json:"endTime"`
}

type Owner struct {
	ID					string `json:"KGId"`
	CampaignID			string `json:"campaignId"`
	Envelope			string `json:"envelope"`
	PrivacyPreference	string `json:"privacyPreference"`
}

type AnonymizedKG struct {
	ID					string `json:"id"`
	CampaignID			string `json:"campaignId"`
	Envelope			string `json:"envelope"`
	Signature			string `json:"signature"`
}

func (s *SmartContract) Test(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, startTime string, endTime string) error {
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

func (s *SmartContract) RetrieveEnvelope(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	ownerBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "" , fmt.Errorf("Failed to read envelope %s from world state. %v", id, err)
	}

	fmt.Print(ownerBytes)

	// var owner Owner
	// error := json.Unmarshal(ownerBytes, owner)
	// if error != nil {
	// 	return "", fmt.Errorf("Failed to unmarshal envelope %v", error)
	// }

	return string(ownerBytes), nil
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

func (s *SmartContract) ShareKnowledgeGraph(ctx contractapi.TransactionContextInterface, id, campaignId, privacyPreference, envelope string) error {
	dataBytes, errId := ctx.GetStub().GetState(id)
	if errId != nil {
		return fmt.Errorf("Failed to read from world state %s %v", id, errId)
	}
	if dataBytes != nil {
		return fmt.Errorf("Data ID %s already existent", id)
	}

	exists, err := s.CampaignExists(ctx, campaignId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Campaign %s does not exist", campaignId)
	}

	owner := Owner{			
		ID:					id,
		CampaignID: 		campaignId,
		Envelope: 			envelope,	
		PrivacyPreference:	privacyPreference,
	}

	ownerJSON, err := json.Marshal(owner)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, ownerJSON)

	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) ShareKGForVerification(ctx contractapi.TransactionContextInterface, id, campaignId, envelope, signature string) error {
	dataBytes, errId := ctx.GetStub().GetState(id)
	if errId != nil {
		return fmt.Errorf("Failed to read from world state %s %v", id, errId)
	}
	if dataBytes != nil {
		return fmt.Errorf("Data ID %s already existent", id)
	}

	exists, err := s.CampaignExists(ctx, campaignId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Campaign %s does not exist", campaignId)
	}

	anonymizedKG := AnonymizedKG{			
		ID:					id,
		CampaignID: 		campaignId,
		Envelope: 			envelope,	
		Signature:			signature,
	}

	ownerJSON, err := json.Marshal(anonymizedKG)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, ownerJSON)

	if err != nil {
		return err
	}

	return nil
}

//TODO dummy function for now
func (s *SmartContract) VerifyProof(ctx contractapi.TransactionContextInterface, userCommit, rollupCommit string) bool {
	return userCommit == rollupCommit
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
