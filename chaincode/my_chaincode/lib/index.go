package lib

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// ------------------------------------

type Admin_Super struct {
	Password string `json:"password"`
	UID      string `json:"uid"`  //  Inspector General of Registration
	Name     string `json:"name"` //  Inspector General of Registration
}

// OfficeCode: Tri letter unique code give to each Sub-Registrar's office
type Admin_OfficeCode struct {
	Password  string   `json:"password"`
	UID       string   `json:"uid"`       // Sub-Registrar
	Name      string   `json:"name"`      // Sub-Registrar
	ToApprove []string `json:"toApprove"` // transactions to approve
}

type User struct {
	Password  string          `json:"password"`
	UID       string          `json:"uid"`
	Name      string          `json:"name"`
	Status    int             `json:"status"` // 0/1/2 - Not verified/Verified/Suspended
	Owned     []string        `json:"owned"`
	Requested []Request_Buyer `json:"requested"`
}

type Request struct {
	Buyer         string    `json:"buyer"`
	Name          string    `json:"name"`
	ProposedPrice int       `json:"proposedPrice"`
	DateTime      time.Time `json:"dateTime"`
}

type Request_Buyer struct {
	ULPIN         string    `json:"ulpin"`
	ProposedPrice int       `json:"proposedPrice"`
	DateTime      time.Time `json:"dateTime"`
}

type Transaction struct {
	Seller              string    `json:"seller"`
	Buyer               string    `json:"buyer"`
	TransactionDateTime time.Time `json:"transactionDateTime"` // when seller/owner accepted the request
	OfficeCode          string    `json:"officeCode"`          // Where estate resides
	ApprovedBy          string    `json:"approvedBy"`          // uid
	ApprovedDateTime    time.Time `json:"approvedDateTime"`
	Price               int       `json:"price"`  // accepted buy seller/owner
	Reason              string    `json:"reason"` // sell, inheritance, gift
}

type Estate struct {
	Owner             string    `json:"owner"`             // uid
	OfficeCode        string    `json:"officeCode"`        // Where estate resides
	Location          string    `json:"location"`          // address
	Area              int       `json:"area"`              // in sq mtr
	Status            int       `json:"status"`            // 0/1/2 - Not verified/Verified/Suspended
	PurchasedOn       time.Time `json:"purchasedOn"`       // current owner since
	SaleAvailability  bool      `json:"saleAvailability"`  // bool
	TransactionsCount int       `json:"transactionsCount"` // total transactions till now
	Requests          []Request `json:"requests"`          // all request from buyers
	BeingSold         bool      `json:"beingSold"`         // true when a request from buyer is accepted
}

// struct for events
/*
type Transaction_Event struct {
	ULPIN         string `json:"ulpin"`
	TransactionCount int    `json:"transactionCount"`
	Transaction      Transaction
}
*/

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
	fmt.Print("\n")
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

// remove this function in future and use builtin methods

func searchArray(arr []string, val string) int {
	for i, s := range arr {
		if s == val {
			return i
		}
	}
	return -1
}
