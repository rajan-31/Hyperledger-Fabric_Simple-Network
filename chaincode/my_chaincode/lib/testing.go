package lib

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Testing

func (s *SmartContract) GetValue(ctx contractapi.TransactionContextInterface, _key string) (map[string]interface{}, error) {

	// this is to handle the data from unknown/misc structs
	data := make(map[string]interface{})

	dataAsBytes, err0 := ctx.GetStub().GetState(_key)

	if err0 != nil {
		return data, fmt.Errorf("GetValue >> Failed to read from world state. %s", err0.Error())
	}

	if dataAsBytes == nil {
		return data, fmt.Errorf("GetValue >> %s does not exist", _key)
	}

	err1 := json.Unmarshal(dataAsBytes, &data)
	if err1 != nil {
		return data, fmt.Errorf("GetValue >> Can't Unmarshal Data")
	}

	fmt.Println(data)
	fmt.Print("\n")

	return data, nil
}

func (s *SmartContract) DeleteValue(ctx contractapi.TransactionContextInterface, _key string) error {

	err0 := ctx.GetStub().DelState(_key)
	if err0 != nil {
		return fmt.Errorf("DeleteValue >> Can't Delete Value for %s. %s", _key, err0.Error())
	}

	fmt.Println("Deleted value for key:", _key)
	fmt.Print("\n")

	return nil
}

func (s *SmartContract) GetAll(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]string, error) {

	resultsIterator, err0 := ctx.GetStub().GetStateByRange(startKey, endKey)

	arrMap := []string{}

	if err0 != nil {
		return arrMap, err0
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err1 := resultsIterator.Next()

		if err1 != nil {
			return arrMap, err1
		}

		// data := make(map[string]interface{})
		// err2 := json.Unmarshal(queryResponse.Value, &data)
		// if err2 != nil {
		// 	return arrMap, fmt.Errorf("GetAll >> Can't Unmarshal Data")
		// }

		data := string(queryResponse.Value)

		fmt.Println("=====================================")
		fmt.Println("Key: "+queryResponse.Key+", Value: ", data)
		fmt.Println("=====================================")

		arrMap = append(arrMap, "Key: "+string(queryResponse.Key)+", Value: "+data)
	}

	fmt.Print("\n")

	return arrMap, nil
}

func (s *SmartContract) DeleteAll(ctx contractapi.TransactionContextInterface, startKey string, endKey string) error {

	resultsIterator, err0 := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err0 != nil {
		return err0
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err1 := resultsIterator.Next()

		if err1 != nil {
			return err1
		}

		key := queryResponse.Key

		err2 := ctx.GetStub().DelState(key)

		if err2 != nil {
			return err2
		}
	}

	fmt.Print("Deleted all/multple values")

	return nil
}
