package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"my_chaincode/lib"
)

func main() {

	chaincode, err := contractapi.NewChaincode(new(lib.SmartContract))

	if err != nil {
		fmt.Printf("Error create Real Estate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Real Estate chaincode: %s", err.Error())
	}

}

// ------------------------------------
//=====================================
