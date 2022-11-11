package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type OwnerData struct {
	Id					string 	`json:"id"`
	AssetType			string	`json:"assetType"`
	CampaignId			string 	`json:"campaignId"`
	Envelope			string 	`json:"envelope"`
	PrivacyPreference	string 	`json:"privacyPreference"`
	Url					string 	`json:"url"`
}


func (s *SmartContract) ShareData(ctx contractapi.TransactionContextInterface, id, campaignId, envelope, privacyPreference string) error {
    idExists, err := s.dataExists(ctx, id)
    if err != nil {
        return err
    }
    if idExists {
        return fmt.Errorf("Id %s already exists", id)
    }

    exists, err := s.campaignExists(ctx, campaignId)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Campaign %s does not exist", campaignId)
    }

    owner := OwnerData{         
        Id:             	id,
		AssetType:	 		string(OnwerDataAssetType),
		CampaignId: campaignId,
        Envelope:           envelope,   
        PrivacyPreference:  privacyPreference,
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

func (s *SmartContract) DeleteSharedData(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.dataExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Error while deleting data: the data %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
}

func (s *SmartContract) dataExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	dataBytes, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("Failed to read data %s from world state. %v", id, err)
    }
	if dataBytes == nil {
		return false, nil
	}

	dataString := string(dataBytes)

	if !(strings.Contains(dataString, string(OnwerDataAssetType))) {
		return false, nil
	}

    return true, nil
}