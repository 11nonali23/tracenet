package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type AssetType string

const (
	CampaignAssetType AssetType 	= "campaign"
	OnwerDataAssetType        		= "ownerData"
	AnonymizedKGAssetType        	= "AnonymizedKG"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
