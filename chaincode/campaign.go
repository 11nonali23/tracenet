package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Campaign struct {
	Id        		string 			`json:"id"`
	AssetType		string			`json:"assetType"`
	Name      		string 			`json:"name"`
	StartTime 		string 			`json:"startTime"`
	EndTime			string 			`json:"endTime"`
}

func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, startTime string, endTime string) error {
	campaign := Campaign{
		Id:       			id,
		AssetType: 			string(CampaignAssetType),
		Name:      			name,
		StartTime: 			startTime,
		EndTime:   			endTime,
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
    exists, err := s.campaignExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Error while deleting campaign: the campaign %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
}

func (s *SmartContract) campaignExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	campaignBytes, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("Failed to read campaign %s from world state. %v", id, err)
    }
	if campaignBytes == nil {
		return false, nil
	}

	campaignString := string(campaignBytes)

	if !(strings.Contains(campaignString, string(CampaignAssetType))) {
		return false, nil
	}

	return true, nil
}







func (s *SmartContract) QueryCampaign(ctx contractapi.TransactionContextInterface, id string) Campaign {
	campaign, _ := s.getCampaign(ctx, id)

	return *campaign
}

func (s *SmartContract) getCampaign(ctx contractapi.TransactionContextInterface, id string) (*Campaign, error) {
	campaignBytes, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("Failed to read campaign %s from world state. %v", id, err)
    }
	if campaignBytes == nil {
		return nil, fmt.Errorf("Campaign %s does not exist", id)
	}

	campaign := new(Campaign)
	_ = json.Unmarshal(campaignBytes, campaign)

    return campaign, nil
}