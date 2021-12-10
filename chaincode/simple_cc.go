package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// ------------------------------------

type DateTime struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

type Admin_Super struct {
	Password string `json:"password"`
	Name     string `json:"name"` //  Inspector General of Registration
	UID      string `json:"uid"`  //  Inspector General of Registration
}

// OfficeCode: Tri letter unique code give to each Sub-Registrar's office
type Admin_OfficeCode struct {
	Password string `json:"password"`
	UID      string `json:"uid"` // Sub-Registrar
}

type User struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	UID      string `json:"uid"`
	Status   int    `json:"status"` // 0/1/2 - Not verified/Verified/Suspended
}

type Estate struct {
	Owner             string    `json:"owner"`      // uid
	OfficeCode        string    `json:"officeCode"` // Where estate resides
	Location          string    `json:"location"`   // address
	Area              int       `json:"area"`       // in sq mtr
	Status            int       `json:"status"`     // 0/1/2 - Not verified/Verified/Suspended
	PurchasedOn       *DateTime //					// current owner since
	SaleAvailability  bool      `json:"saleAvailability"`  // bool
	TransactionsCount int       `json:"transactionsCount"` // total transactions till now
}

// ------------------------------------

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	// put super user name and password
	key := "admin_super"
	data := Admin_Super{
		Password: "123456",
		Name:     "FName MName LName",
		UID:      "123456789012",
	}

	marshaled_data, _ := json.Marshal(data)

	err := ctx.GetStub().PutState(key, marshaled_data)

	if err != nil {
		return fmt.Errorf("InitLedger >> failed to put to world state. %s", err.Error())
	}

	fmt.Println("=====================================")
	fmt.Println("Chaincode Initated.")
	fmt.Println("=====================================")

	return nil
}

// ------------------------------------

// Helper Functions - Private

func (s *SmartContract) verifyPassword(ctx contractapi.TransactionContextInterface, _username string, _password string) (bool, error) {

	// get data
	dataAsBytes, err0 := ctx.GetStub().GetState(_username)

	if err0 != nil {
		return false, fmt.Errorf("GetPassword >> Failed to read from world state. %s", err0.Error())
	}

	if dataAsBytes == nil {
		return false, fmt.Errorf("GetPassword >> %s does not exist", _username)
	}

	// extract password
	// this is to handle the data from unknown/misc structs
	data := make(map[string]interface{})
	err1 := json.Unmarshal(dataAsBytes, &data)
	if err1 != nil {
		return false, fmt.Errorf("GetPassword >> Can't Unmarshal Data")
	}

	password, ok := data["password"].(string)
	if !ok {
		// password is not a string
		return false, fmt.Errorf("GetPassword >> Password is not a string")
	}

	if _password == password {
		return true, nil
	}

	return false, nil
}

// ------------------------------------

// For Admin super

func (s *SmartContract) CreateOrModify_Admin(ctx contractapi.TransactionContextInterface, _password string, officeCode string, newAdminPassword string, uid string) error {

	verified, err0 := s.verifyPassword(ctx, "admin_super", _password)

	if err0 != nil {
		return fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return fmt.Errorf("CreateOrModify_Admin >> Password Missmatched for %s", "admin_super")
	}

	//=====================================

	key := "admin_" + officeCode
	data := Admin_OfficeCode{
		Password: newAdminPassword,
		UID:      uid,
	}

	marshaled_data, _ := json.Marshal(data)
	err1 := ctx.GetStub().PutState(key, marshaled_data)
	if err1 != nil {
		return fmt.Errorf("CreateOrModify_Admin >> Failed to put to world state. %s", err1.Error())
	}
	return nil
}

// For Admin

func (s *SmartContract) Create_User(ctx contractapi.TransactionContextInterface, _username string, _password string, uid string, newUserPassword string, name string) (User, error) {
	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return User{}, fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return User{}, fmt.Errorf("Create_User >> Password Missmatched for %s", _username)
	}

	//=====================================

	key := "user" + "_" + uid
	data := User{
		Password: newUserPassword,
		Name:     name,
		UID:      uid,
		Status:   0,
	}

	marshaled_data, _ := json.Marshal(data)
	err1 := ctx.GetStub().PutState(key, marshaled_data)
	if err1 != nil {
		return User{}, fmt.Errorf("Create_User >> Failed to put to world state. %s", err1.Error())
	}
	return data, nil
}

func (s *SmartContract) Create_Estate(ctx contractapi.TransactionContextInterface, _username string, _password string, serveyNo string, owner string, location string, area int, purchasedOn DateTime, transactionsCount int) (Estate, error) {
	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return Estate{}, fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return Estate{}, fmt.Errorf("Create_Estate >> Password Missmatched for %s", _username)
	}

	//=====================================

	key := "estate" + "_" + serveyNo
	data := Estate{
		Owner:             owner,
		OfficeCode:        strings.TrimPrefix(_username, "admin_"),
		Location:          location,
		Area:              area,
		Status:            0,
		PurchasedOn:       &purchasedOn,
		SaleAvailability:  false,
		TransactionsCount: transactionsCount,
	}

	marshaled_data, _ := json.Marshal(data)
	err1 := ctx.GetStub().PutState(key, marshaled_data)
	if err1 != nil {
		return Estate{}, fmt.Errorf("Create_Estate >> Failed to put to world state. %s", err1.Error())
	}
	return data, nil
}

// For System

func (s *SmartContract) Verify_User(ctx contractapi.TransactionContextInterface, _username string, _password string, uid string, status int) error {
	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return fmt.Errorf("Verify_User >> Password Missmatched for %s", _username)
	}

	//=====================================

	key := "user" + "_" + uid
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

	user.Status = status

	//=====================================

	marshaled_data, _ := json.Marshal(user)
	err3 := ctx.GetStub().PutState(key, marshaled_data)
	if err3 != nil {
		return fmt.Errorf("Verify_User >> Failed to put to world state. %s", err3.Error())
	}

	return nil
}

// ------------------------------------

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

	return data, nil
}

func (s *SmartContract) GetAll(ctx contractapi.TransactionContextInterface) error {

	startKey := ""
	endKey := ""

	resultsIterator, err0 := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err0 != nil {
		return err0
	}
	defer resultsIterator.Close()

	// count := 0
	// for resultsIterator.HasNext() {
	// 	count++
	// }
	// allValues := make([]map[string]interface{}, count)

	for resultsIterator.HasNext() {
		queryResponse, err1 := resultsIterator.Next()

		if err1 != nil {
			return err1
		}

		data := make(map[string]interface{})
		err2 := json.Unmarshal(queryResponse.Value, &data)
		if err2 != nil {
			return fmt.Errorf("GetAll >> Can't Unmarshal Data")
		}

		fmt.Println("=====================================")
		fmt.Println("Key: "+queryResponse.Key+", Value: ", data)
		fmt.Println("=====================================")

		// data["Key"] = queryResponse.Key
		// data["Values"] ueryResponse.Value
		// allValues = append(allValues, data)
	}

	// println(allValues)

	return nil
}

// ------------------------------------

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create Real Estate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Real Estate chaincode: %s", err.Error())
	}
}

//=====================================

// ------------------------------------
