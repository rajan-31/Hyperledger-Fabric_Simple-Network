package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// For Admin super

func (s *SmartContract) CreateOrModify_Admin(ctx contractapi.TransactionContextInterface, _username string, _password string, officeCode string, newAdminPassword string, uid string, name string) error {

	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return fmt.Errorf("CreateOrModify_Admin >> Password Missmatched for %s", _username)
	}

	//=====================================

	key := "admin_" + officeCode
	data := Admin_OfficeCode{
		Password:  newAdminPassword,
		UID:       uid,
		Name:      name,
		ToApprove: []string{},
	}

	marshaled_data, _ := json.Marshal(data)
	err1 := ctx.GetStub().PutState(key, marshaled_data)
	if err1 != nil {
		return fmt.Errorf("CreateOrModify_Admin >> Failed to put to world state. %s", err1.Error())
	}
	return nil
}

// For Admin

func (s *SmartContract) Modify_User(ctx contractapi.TransactionContextInterface, uid string, name string, status int) (User, error) {

	key := "user" + "_" + uid
	user := new(User)

	// get user data
	dataAsBytes, err1 := ctx.GetStub().GetState(key)

	if err1 != nil {
		return User{}, fmt.Errorf("Modify_User >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return User{}, fmt.Errorf("Modify_User >> %s does not exist", key)
	}

	err2 := json.Unmarshal(dataAsBytes, &user)
	if err2 != nil {
		return User{}, fmt.Errorf("Modify_User >> Can't Unmarshal Data")
	}

	//=====================================
	temp_status := user.Status
	if status != -1 {
		temp_status = status
	}
	data := User{
		Password:  user.Password,
		Name:      name,
		UID:       user.UID,
		Status:    temp_status,
		Owned:     user.Owned,
		Requested: user.Requested,
	}

	marshaled_data, _ := json.Marshal(data)
	err3 := ctx.GetStub().PutState(key, marshaled_data)
	if err3 != nil {
		return User{}, fmt.Errorf("Modify_User >> Failed to put to world state. %s", err3.Error())
	}

	return data, nil
}

func (s *SmartContract) Create_Estate(ctx contractapi.TransactionContextInterface, officeCode string, ulpin string, owner string, location string, area int, purchasedOn string, transactionsCount int) (Estate, error) {

	key := "estate" + "_" + ulpin
	temp_dateTime, _ := time.Parse(time.RFC3339, purchasedOn) // purchasedOn => 2021-12-15T20:34:33+05:30
	data := Estate{
		Owner:             owner,
		OfficeCode:        officeCode,
		Location:          location,
		Area:              area,
		Status:            0,
		PurchasedOn:       temp_dateTime,
		SaleAvailability:  false,
		TransactionsCount: transactionsCount,
		Requests:          []Request{},
		BeingSold:         false,
	}

	marshaled_data, _ := json.Marshal(data)
	err1 := ctx.GetStub().PutState(key, marshaled_data)
	if err1 != nil {
		return Estate{}, fmt.Errorf("Create_Estate >> Failed to put to world state. %s", err1.Error())
	}

	//=====================================

	// get data
	key = "user" + "_" + owner
	dataAsBytes, err2 := ctx.GetStub().GetState(key)

	if err2 != nil {
		return data, fmt.Errorf("Create_Estate >> Failed to read from world state. %s", err2.Error())
	}

	if dataAsBytes == nil {
		return data, fmt.Errorf("Create_Estate >> %s does not exist", key)
	}

	user := new(User)
	err3 := json.Unmarshal(dataAsBytes, &user)
	if err3 != nil {
		return data, fmt.Errorf("GetValue >> Can't Unmarshal Data")
	}

	temp_owned := user.Owned
	i := searchArray(temp_owned, ulpin)
	if i != -1 {
		return data, fmt.Errorf("Create_Estate >> User alredy own estate with ulpin: %s", ulpin)
	}

	temp_owned = append(temp_owned, ulpin)
	user.Owned = temp_owned

	key2 := "user" + "_" + owner
	marshaled_data2, _ := json.Marshal(user)
	err4 := ctx.GetStub().PutState(key2, marshaled_data2)
	if err4 != nil {
		return Estate{}, fmt.Errorf("Create_Estate >> Failed to put to world state. %s", err4.Error())
	}

	return data, nil
}

func (s *SmartContract) Modify_Estate(ctx contractapi.TransactionContextInterface, officeCode string, ulpin string, location string, area int, purchasedOn string, transactionsCount int) (Estate, error) {

	// get data
	key := "estate" + "_" + ulpin
	dataAsBytes, err1 := ctx.GetStub().GetState(key)

	if err1 != nil {
		return Estate{}, fmt.Errorf("Modify_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return Estate{}, fmt.Errorf("Modify_Estate >> %s does not exist", key)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes, &estate)
	if err2 != nil {
		return Estate{}, fmt.Errorf("GetValue >> Can't Unmarshal Data")
	}

	//=====================================

	var temp_dateTime time.Time
	if location == "" {
		location = estate.Location
	}
	if area == -1 {
		area = estate.Area
	}
	if purchasedOn == "" {
		temp_dateTime = estate.PurchasedOn
	} else {
		temp_dateTime, _ = time.Parse(time.RFC3339, purchasedOn) // purchasedOn => 2021-12-15T20:34:33+05:30
	}
	if transactionsCount == -1 {
		transactionsCount = estate.TransactionsCount
	}

	data := Estate{
		Owner:             estate.Owner,
		OfficeCode:        officeCode,
		Location:          location,
		Area:              area,
		Status:            estate.Status,
		PurchasedOn:       temp_dateTime,
		SaleAvailability:  estate.SaleAvailability,
		TransactionsCount: transactionsCount,
		Requests:          estate.Requests,
		BeingSold:         estate.BeingSold,
	}

	marshaled_data, _ := json.Marshal(data)
	err3 := ctx.GetStub().PutState(key, marshaled_data)
	if err3 != nil {
		return Estate{}, fmt.Errorf("Modify_Estate >> Failed to put to world state. %s", err3.Error())
	}

	return data, nil
}

func (s *SmartContract) Add_Transaction(ctx contractapi.TransactionContextInterface, ulpin string, num int, seller string, buyer string, reason string, proposedPrice int, tDateTime string, officeCode string, approvedBy string, aDateTime string) (Transaction, error) {
	key := "transaction" + "_" + ulpin + "_" + strconv.Itoa(num)

	temp_tDateTime, _ := time.Parse(time.RFC3339, tDateTime)
	temp_aDateTime, _ := time.Parse(time.RFC3339, aDateTime)
	data := Transaction{
		Seller:              seller,
		Buyer:               buyer,
		TransactionDateTime: temp_tDateTime,
		OfficeCode:          officeCode,
		ApprovedBy:          approvedBy,
		ApprovedDateTime:    temp_aDateTime,
		Price:               proposedPrice,
		Reason:              reason,
	}

	marshaled_data0, _ := json.Marshal(data)
	err0 := ctx.GetStub().PutState(key, marshaled_data0)
	if err0 != nil {
		return Transaction{}, fmt.Errorf("Add_Transaction >> failed to put to world state. %s", err0.Error())
	}

	return data, nil
}

func (s *SmartContract) ApproveSell_Estate(ctx contractapi.TransactionContextInterface, _username string, ulpin string, dateTime string) (Estate, error) {

	// get data of estate
	key1 := "estate" + "_" + ulpin
	dataAsBytes0, err1 := ctx.GetStub().GetState(key1)

	if err1 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes0 == nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> %s does not exist", key1)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes0, &estate)
	if err2 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Can't Unmarshal Data")
	}

	// get data of transaction
	key2 := "transaction" + "_" + ulpin + "_" + strconv.Itoa(estate.TransactionsCount+1)
	dataAsBytes1, err2 := ctx.GetStub().GetState(key2)

	if err2 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Failed to read from world state. %s", err2.Error())
	}

	if dataAsBytes1 == nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> %s does not exist", key2)
	}

	transaction := new(Transaction)
	err3 := json.Unmarshal(dataAsBytes1, &transaction)
	if err3 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Can't Unmarshal Data")
	}

	// update approvedBy in transaction
	key3 := _username
	dataAsBytes2, err3 := ctx.GetStub().GetState(key3)

	if err3 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Failed to read from world state. %s", err3.Error())
	}

	if dataAsBytes1 == nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> %s does not exist", key3)
	}

	admin := new(Admin_OfficeCode)
	err4 := json.Unmarshal(dataAsBytes2, &admin)
	if err4 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Can't Unmarshal Data")
	}

	temp_dateTime, _ := time.Parse(time.RFC3339, dateTime)

	transaction.ApprovedBy = admin.UID
	transaction.ApprovedDateTime = temp_dateTime

	// update owner, purchasedOn of estate
	estate.Owner = transaction.Buyer
	estate.PurchasedOn = temp_dateTime
	estate.SaleAvailability = false
	estate.TransactionsCount++
	estate.BeingSold = false

	// update transaction
	marshaled_data0, _ := json.Marshal(estate)
	err5 := ctx.GetStub().PutState(key1, marshaled_data0)
	if err5 != nil {
		return Estate{}, fmt.Errorf("ApproveSell_Estate >> Failed to put to world state. %s", err5.Error())
	}

	// update estate
	marshaled_data1, _ := json.Marshal(transaction)
	err6 := ctx.GetStub().PutState(key2, marshaled_data1)
	if err6 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Failed to put to world state. %s", err6.Error())
	}

	//=====================================

	// remove estate from seller's owned

	key4 := "user" + "_" + transaction.Seller
	user0 := new(User)

	// get user0 data
	dataAsBytes3, err7 := ctx.GetStub().GetState(key4)

	if err7 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Failed to read from world state. %s", err7.Error())
	}

	if dataAsBytes3 == nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> %s does not exist", key4)
	}

	err8 := json.Unmarshal(dataAsBytes3, &user0)
	if err8 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Can't Unmarshal Data")
	}

	temp_owned0 := user0.Owned
	i0 := searchArray(temp_owned0, ulpin)
	user0.Owned = append(temp_owned0[:i0], temp_owned0[i0+1:]...)

	marshaled_data2, _ := json.Marshal(user0)
	err9 := ctx.GetStub().PutState(key4, marshaled_data2)
	if err9 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Failed to put to world state. %s", err9.Error())
	}

	// add estate to buyer's owned

	key5 := "user" + "_" + transaction.Buyer
	user1 := new(User)

	// get user1 data
	dataAsBytes4, err10 := ctx.GetStub().GetState(key5)

	if err10 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Failed to read from world state. %s", err10.Error())
	}

	if dataAsBytes4 == nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> %s does not exist", key5)
	}

	err11 := json.Unmarshal(dataAsBytes4, &user1)
	if err11 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Can't Unmarshal Data")
	}

	temp_owned1 := user1.Owned
	user1.Owned = append(temp_owned1, ulpin)

	marshaled_data3, _ := json.Marshal(user1)
	err12 := ctx.GetStub().PutState(key5, marshaled_data3)
	if err12 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> Failed to put to world state. %s", err12.Error())
	}

	//=====================================
	// remove from admin toApprove

	temp_toApprove := admin.ToApprove
	i1 := searchArray(temp_toApprove, key2)
	admin.ToApprove = append(temp_toApprove[:i1], temp_toApprove[i1+1:]...)

	marshaled_data4, _ := json.Marshal(admin)
	err13 := ctx.GetStub().PutState(key3, marshaled_data4)
	if err13 != nil {
		return *estate, fmt.Errorf("ApproveSell_Estate >> failed to put to world state. %s", err13.Error())
	}

	return *estate, nil
}

func (s *SmartContract) RejectSell_Estate(ctx contractapi.TransactionContextInterface, _username string, ulpin string) (Estate, error) {

	// get estate data

	key1 := "estate" + "_" + ulpin
	dataAsBytes1, err1 := ctx.GetStub().GetState(key1)

	if err1 != nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes1 == nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> %s does not exist", key1)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes1, &estate)
	if err2 != nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> Can't Unmarshal Data")
	}

	// get transaction data

	key2 := "transaction" + "_" + ulpin + "_" + strconv.Itoa(estate.TransactionsCount+1)
	dataAsBytes2, err3 := ctx.GetStub().GetState(key2)

	if err3 != nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> Failed to read from world state. %s", err3.Error())
	}

	if dataAsBytes2 == nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> %s does not exist", key2)
	}

	transaction := new(Transaction)
	err4 := json.Unmarshal(dataAsBytes2, &transaction)
	if err4 != nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> Can't Unmarshal Data")
	}

	// get admin data

	key3 := "admin" + "_" + estate.OfficeCode
	dataAsBytes3, err5 := ctx.GetStub().GetState(key3)

	if err5 != nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> Failed to read from world state. %s", err5.Error())
	}

	if dataAsBytes3 == nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> %s does not exist", key3)
	}

	admin := new(Admin_OfficeCode)
	err6 := json.Unmarshal(dataAsBytes3, &admin)
	if err6 != nil {
		return Estate{}, fmt.Errorf("RejectSell_Estate >> Can't Unmarshal Data")
	}

	// change beingSold

	estate.BeingSold = false

	marshaled_data1, _ := json.Marshal(estate)
	err7 := ctx.GetStub().PutState(key1, marshaled_data1)
	if err7 != nil {
		return *estate, fmt.Errorf("RejectSell_Estate >> failed to put to world state. %s", err7.Error())
	}

	// delete transaction

	err8 := ctx.GetStub().DelState(key2)
	if err8 != nil {
		return *estate, fmt.Errorf("RejectSell_Estate >> failed to delete from world state. %s", err8.Error())
	}

	// remove from admin toApprove

	temp_toApprove := admin.ToApprove
	i := searchArray(temp_toApprove, key2)
	admin.ToApprove = append(temp_toApprove[:i], temp_toApprove[i+1:]...)

	marshaled_data2, _ := json.Marshal(admin)
	err9 := ctx.GetStub().PutState(key3, marshaled_data2)
	if err9 != nil {
		return *estate, fmt.Errorf("RejectSell_Estate >> failed to put to world state. %s", err9.Error())
	}

	return *estate, nil

}
