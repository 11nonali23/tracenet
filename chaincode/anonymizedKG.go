package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AnonymizedKG struct {
	Id					string 	`json:"id"`
	AssetType			string	`json:"assetType"`
	CampaignId			string 	`json:"campaignId"`
	RecipientId			string 	`json:"recipientId"`
	RollupEnvelope		string 	`json:"rollupEnvelope"`
	RecipientEnvelope	string 	`json:"recipientEnvelope"`
	Signature			string 	`json:"signature"`
}

func (s *SmartContract) ShareAnonymizedKGForVerification(ctx contractapi.TransactionContextInterface, id, campaignId, recipientId, rollupEnvelope, signature string) error {
    idExists, err := s.anonymizedKGExists(ctx, id)
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

    anonymizedKG := AnonymizedKG{         
        Id:             	id,
		AssetType:	 		string(OnwerDataAssetType),
		CampaignId: 		campaignId,
		RecipientId: 		recipientId,	
        RollupEnvelope:     rollupEnvelope,
		RecipientEnvelope: 	"",   
        Signature:  		signature,
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

func (s *SmartContract) ShareAnonymizedKGWithRecipient(ctx contractapi.TransactionContextInterface, id, campaignId, recipientId, recipientEnvelope string) error {
    idExists, err := s.anonymizedKGExists(ctx, id)
    if err != nil {
        return err
    }
    if !idExists {
        return fmt.Errorf("Id %s does not exist", id)
    }

    exists, err := s.campaignExists(ctx, campaignId)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Campaign %s does not exist", campaignId)
    }

    anonymizedKG, err := s.getAnonymizedKG(ctx, id)
	if err != nil {
		return fmt.Errorf("Campaign %s does not exist", campaignId)
	}
	anonymizedKG.RecipientEnvelope = recipientEnvelope

    anonymizedKGJSON, err := json.Marshal(anonymizedKG)
    if err != nil {
        return err
    }

	return ctx.GetStub().PutState(id, anonymizedKGJSON)
}

func (s *SmartContract) DeleteAnonymizedKG(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.anonymizedKGExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Error while deleting data: the data %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
}

func (s *SmartContract) anonymizedKGExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	KGBytes, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("Failed to read KG %s from world state. %v", id, err)
    }
	if KGBytes == nil {
		return false, nil
	}

	KGString := string(KGBytes)

	if !(strings.Contains(KGString, string(AnonymizedKGAssetType))) {
		return false, nil
	}

    return true, nil
}

func (s *SmartContract) getAnonymizedKG(ctx contractapi.TransactionContextInterface, id string) (*AnonymizedKG, error) {
	anonymizedKGBytes, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("Failed to read anonymized KG %s from world state. %v", id, err)
    }
	if anonymizedKGBytes == nil {
		return nil, fmt.Errorf("anonymized KG %s does not exist", id)
	}

	anonymizedKG := new(AnonymizedKG)
	_ = json.Unmarshal(anonymizedKGBytes, anonymizedKG)

    return anonymizedKG, nil
}