package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// User

func (s *SmartContract) ChangeAvail_Estate(ctx contractapi.TransactionContextInterface, ulpin string, saleAvailability bool) error {
	// get data
	key := "estate" + "_" + ulpin
	dataAsBytes, err1 := ctx.GetStub().GetState(key)

	if err1 != nil {
		return fmt.Errorf("ChangeAvail_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return fmt.Errorf("ChangeAvail_Estate >> %s does not exist", key)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes, &estate)
	if err2 != nil {
		return fmt.Errorf("GetValue >> Can't Unmarshal Data")
	}

	//=====================================

	estate.SaleAvailability = saleAvailability

	marshaled_data, _ := json.Marshal(estate)
	err3 := ctx.GetStub().PutState(key, marshaled_data)
	if err3 != nil {
		return fmt.Errorf("ChangeAvail_Estate >> Failed to put to world state. %s", err3.Error())
	}

	return nil
}

func (s *SmartContract) RequestToBuy_Estate(ctx contractapi.TransactionContextInterface, _buyer string, _name string, ulpin string, proposedPrice int, dateTime string) (Request, error) {
	// get data estate
	key := "estate" + "_" + ulpin
	dataAsBytes, err1 := ctx.GetStub().GetState(key)

	if err1 != nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> %s does not exist", ulpin)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes, &estate)
	if err2 != nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> Can't Unmarshal Data")
	}

	// get data buyer

	key2 := "user" + "_" + _buyer
	dataAsBytes2, err3 := ctx.GetStub().GetState(key2)

	if err3 != nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> Failed to read from world state. %s", err3.Error())
	}

	if dataAsBytes2 == nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> %s does not exist", key2)
	}

	buyer := new(User)
	err4 := json.Unmarshal(dataAsBytes2, &buyer)
	if err4 != nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> Can't Unmarshal Data")
	}

	//=====================================

	// add or update request in estate array

	temp_requests := estate.Requests
	temp_dateTime, _ := time.Parse(time.RFC3339, dateTime)

	flag := false
	index := -1
	for i, r := range temp_requests {
		if r.Buyer == _buyer {
			temp_requests[i].ProposedPrice = proposedPrice
			temp_requests[i].DateTime = temp_dateTime
			flag = true
			index = i
			break
		}
	}

	if !flag {
		temp_requests = append(temp_requests, Request{
			Buyer:         _buyer,
			Name:          _name,
			ProposedPrice: proposedPrice,
			DateTime:      temp_dateTime,
		})
		index = len(temp_requests) - 1
	}
	estate.Requests = temp_requests

	// add or update request in buyer requested array

	temp_requested := buyer.Requested

	flag2 := false
	for i2, r2 := range temp_requested {
		if r2.ULPIN == ulpin {
			temp_requested[i2].ProposedPrice = proposedPrice
			temp_requested[i2].DateTime = temp_dateTime
			flag2 = true
			break
		}
	}

	if !flag2 {
		temp_requested = append(temp_requested, Request_Buyer{
			ULPIN:         ulpin,
			ProposedPrice: proposedPrice,
			DateTime:      temp_dateTime,
		})
	}
	buyer.Requested = temp_requested

	//=====================================

	// update estate requests array

	marshaled_data, _ := json.Marshal(estate)
	err5 := ctx.GetStub().PutState(key, marshaled_data)
	if err5 != nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> Failed to put to world state. %s", err5.Error())
	}

	// update byer requested

	marshaled_data2, _ := json.Marshal(buyer)
	err6 := ctx.GetStub().PutState(key2, marshaled_data2)
	if err6 != nil {
		return Request{}, fmt.Errorf("RequestToBuy_Estate >> Failed to put to world state. %s", err6.Error())
	}

	return temp_requests[index], nil
}

func (s *SmartContract) AcceptRequest_Estate(ctx contractapi.TransactionContextInterface, _username string, _password string, ulpin string, buyer string, dateTime string, reason string) (Transaction, error) {
	verified, err0 := s.verifyPassword(ctx, _username, _password)

	if err0 != nil {
		return Transaction{}, fmt.Errorf("verifyPassword >> Verify password %s", err0.Error())
	} else if !verified {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Password Missmatched for %s", _username)
	}

	//=====================================

	// get data
	key1 := "estate" + "_" + ulpin
	dataAsBytes, err1 := ctx.GetStub().GetState(key1)

	if err1 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes == nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> %s does not exist", ulpin)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes, &estate)
	if err2 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Can't Unmarshal Data")
	}

	//=====================================

	// check if being sold already
	if estate.BeingSold {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Estate is already being sold")
	}

	// stop from accepting other requests
	estate.BeingSold = true

	// it is updated in later step

	//=====================================

	flag := false
	index := -1
	for i, r := range estate.Requests {
		if r.Buyer == buyer {
			flag = true
			index = i
		}
	}

	if !flag {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> No request found for given buyer: %s", buyer)
	}

	temp_dateTime, _ := time.Parse(time.RFC3339, dateTime)
	temp_transaction := Transaction{
		Seller:              strings.TrimPrefix(_username, "user_"),
		Buyer:               estate.Requests[index].Buyer,
		TransactionDateTime: temp_dateTime,
		OfficeCode:          estate.OfficeCode,
		ApprovedBy:          "",
		ApprovedDateTime:    time.Time{},
		Price:               estate.Requests[index].ProposedPrice,
		Reason:              reason,
	}

	key2 := "transaction" + "_" + ulpin + "_" + strconv.Itoa(estate.TransactionsCount+1)
	marshaled_data1, _ := json.Marshal(temp_transaction)
	err4 := ctx.GetStub().PutState(key2, marshaled_data1)
	if err4 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> failed to put to world state. %s", err4.Error())
	}

	// update being sold flag and delete all requests
	estate.Requests = []Request{}

	marshaled_data0, _ := json.Marshal(estate)
	err3 := ctx.GetStub().PutState(key1, marshaled_data0)
	if err3 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> failed to put to world state. %s", err3.Error())
	}

	//=====================================
	// add request in toApprove of admin

	key3 := "admin" + "_" + estate.OfficeCode
	dataAsBytes1, err4 := ctx.GetStub().GetState(key3)

	if err4 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Failed to read from world state. %s", err4.Error())
	}

	if dataAsBytes1 == nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> %s does not exist", key3)
	}

	admin := new(Admin_OfficeCode)
	err5 := json.Unmarshal(dataAsBytes1, &admin)
	if err5 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Can't Unmarshal Data")
	}

	admin.ToApprove = append(admin.ToApprove, key2)

	marshaled_data2, _ := json.Marshal(admin)
	err6 := ctx.GetStub().PutState(key3, marshaled_data2)
	if err6 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> failed to put to world state. %s", err6.Error())
	}

	//=====================================
	// remove request from buyer Requested
	// get buyer data

	key4 := "user" + "_" + buyer

	dataAsBytes2, err7 := ctx.GetStub().GetState(key4)

	if err7 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Failed to read from world state. %s", err7.Error())
	}

	if dataAsBytes2 == nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> %s does not exist", key4)
	}

	buyer_data := new(User)
	err8 := json.Unmarshal(dataAsBytes2, &buyer_data)
	if err8 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> Can't Unmarshal Data")
	}

	// delete request from buyer requested

	temp_requested := buyer_data.Requested

	flag2 := false
	index2 := -1
	for i2, r2 := range temp_requested {
		if r2.ULPIN == ulpin {
			flag2 = true
			index2 = i2
		}
	}

	if !flag2 {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> No request found for given seveyNo: %s", ulpin)
	}

	buyer_data.Requested = append(temp_requested[:index2], temp_requested[index2+1:]...)

	// update buyer data

	marshaled_data3, _ := json.Marshal(buyer_data)
	err9 := ctx.GetStub().PutState(key4, marshaled_data3)
	if err9 != nil {
		return Transaction{}, fmt.Errorf("AcceptRequest_Estate >> failed to put to world state. %s", err9.Error())
	}

	//=====================================
	// set event

	/*
		transaction_event := Transaction_Event{
			ULPIN:         ulpin,
			TransactionCount: estate.TransactionsCount + 1,
			Transaction:      temp_transaction,
		}

		marshaled_data3, _ := json.Marshal(transaction_event)

		ctx.GetStub().SetEvent("newEstateSell", marshaled_data3)
	*/

	return temp_transaction, nil
}

func (s *SmartContract) ClearRequests_Estate(ctx contractapi.TransactionContextInterface, ulpin string, buyer string) error {
	// get estate data

	key1 := "estate" + "_" + ulpin
	dataAsBytes1, err1 := ctx.GetStub().GetState(key1)

	if err1 != nil {
		return fmt.Errorf("ClearRequests_Estate >> Failed to read from world state. %s", err1.Error())
	}

	if dataAsBytes1 == nil {
		return fmt.Errorf("ClearRequests_Estate >> %s does not exist", key1)
	}

	estate := new(Estate)
	err2 := json.Unmarshal(dataAsBytes1, &estate)
	if err2 != nil {
		return fmt.Errorf("ClearRequests_Estate >> Can't Unmarshal Data")
	}

	//=====================================

	temp_requests := estate.Requests

	buyers_list := []string{}

	// for estate requets
	// one or all requests?

	if buyer == "" {

		// used in removing requested

		for _, request := range temp_requests {
			buyers_list = append(buyers_list, request.Buyer)
		}

		// remove all requests
		estate.Requests = []Request{}
	} else {

		// used in removing requested

		buyers_list = append(buyers_list, buyer)

		// find request and remove

		flag := false
		index := -1
		for i, r := range temp_requests {
			if r.Buyer == buyer {
				flag = true
				index = i
			}
		}

		if !flag {
			return fmt.Errorf("ClearRequests_Estate >> No request found for given buyer: %s", buyer)
		}

		estate.Requests = append(temp_requests[:index], temp_requests[index+1:]...)

	}

	marshaled_data1, _ := json.Marshal(estate)
	err3 := ctx.GetStub().PutState(key1, marshaled_data1)
	if err3 != nil {
		return fmt.Errorf("ClearRequests_Estate >> failed to put to world state. %s", err3.Error())
	}

	//=====================================

	// remove requested from buyer/s

	for _, trav_buyer := range buyers_list {

		// get buyer data

		key2 := "user" + "_" + trav_buyer

		dataAsBytes2, err4 := ctx.GetStub().GetState(key2)

		if err4 != nil {
			return fmt.Errorf("ClearRequests_Estate >> Failed to read from world state. %s", err4.Error())
		}

		if dataAsBytes2 == nil {
			return fmt.Errorf("ClearRequests_Estate >> %s does not exist", key2)
		}

		buyer_data := new(User)
		err5 := json.Unmarshal(dataAsBytes2, &buyer_data)
		if err5 != nil {
			return fmt.Errorf("ClearRequests_Estate >> Can't Unmarshal Data")
		}

		// find and delete requested entry

		temp_requested := buyer_data.Requested

		index2 := -1

		for i0, entry := range temp_requested {
			if entry.ULPIN == ulpin {
				index2 = i0
				break
			}
		}

		buyer_data.Requested = append(temp_requested[:index2], temp_requested[index2+1:]...)

		marshaled_data2, _ := json.Marshal(buyer_data)
		err6 := ctx.GetStub().PutState(key2, marshaled_data2)
		if err6 != nil {
			return fmt.Errorf("ClearRequests_Estate >> failed to put to world state. %s", err6.Error())
		}
	}

	return nil
}
