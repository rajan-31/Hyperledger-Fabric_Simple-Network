package lib

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// For System

func (s *SmartContract) Create_User(ctx contractapi.TransactionContextInterface, uid string, name string) (User, error) {

	key := "user" + "_" + uid
	data := User{
		Password:  uid,
		Name:      name,
		UID:       uid,
		Status:    0,
		Owned:     []string{},
		Requested: []Request_Buyer{},
	}

	marshaled_data, _ := json.Marshal(data)
	err1 := ctx.GetStub().PutState(key, marshaled_data)
	if err1 != nil {
		return User{}, fmt.Errorf("Create_User >> Failed to put to world state. %s", err1.Error())
	}
	return data, nil
}

func (s *SmartContract) Verify_User(ctx contractapi.TransactionContextInterface, _username string, _password string, status int, newPassword string) error {
	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return fmt.Errorf("Verify_User >> Password Missmatched for %s", _username)
	}

	//=====================================

	key := _username
	user := new(User)

	// get user data
	dataAsBytes, err1 := ctx.GetStub().GetState(key)

	if err1 != nil {
		return fmt.Errorf("Verify_User >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return fmt.Errorf("Verify_User >> %s does not exist", key)
	}

	err2 := json.Unmarshal(dataAsBytes, &user)
	if err2 != nil {
		return fmt.Errorf("Verify_User >> Can't Unmarshal Data")
	}

	//=====================================

	user.Status = status
	user.Password = newPassword

	marshaled_data, _ := json.Marshal(user)
	err3 := ctx.GetStub().PutState(key, marshaled_data)
	if err3 != nil {
		return fmt.Errorf("Verify_User >> Failed to put to world state. %s", err3.Error())
	}

	return nil
}

func (s *SmartContract) Verify_Estate(ctx contractapi.TransactionContextInterface, _username string, _password string, ulpin string, status int) error {
	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return fmt.Errorf("Verify_Estate >> Password Missmatched for %s", _username)
	}

	//=====================================

	// get data
	key := "estate" + "_" + ulpin
	dataAsBytes, err1 := ctx.GetStub().GetState(key)

	if err1 != nil {
		return fmt.Errorf("Verify_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return fmt.Errorf("Verify_Estate >> %s does not exist", key)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes, &estate)
	if err2 != nil {
		return fmt.Errorf("GetValue >> Can't Unmarshal Data")
	}

	//=====================================

	estate.Status = status

	marshaled_data, _ := json.Marshal(estate)
	err3 := ctx.GetStub().PutState(key, marshaled_data)
	if err3 != nil {
		return fmt.Errorf("Verify_Estate >> Failed to put to world state. %s", err3.Error())
	}

	return nil
}
