#!/bin/bash

# from mains shell: docker exec cli sh -c 'sh ./scripts/chaincode_dev-main.sh . 5'

# run in peer bash
# 1 = init,     to init chaincode
# 2 = 1/2/3/.../,  commit sequence

_INIT=$1
_SEQ=$2

# if ["$2" = ""];
# then
#     _SEQ=1
# fi


echo -e "\n\n* Change Dir to 'chaincode'"

cd ../../../chaincode



echo -e "\n\n* Set env vars "

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer1.org1.example.com:7051
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem



echo -e "\n\n* Package"

peer lifecycle chaincode package ./packages/simple_cc_v1_v1.tar.gz -p ./my_chaincode --label simplecc_v1_v1



echo -e "\n\n* Install"

peer lifecycle chaincode install ./packages/simple_cc_v1_v1.tar.gz



echo -e "\n\n* Check install"

peer lifecycle chaincode queryinstalled



echo -e "\n\n* Get installed count"

TEMP_COUNT=$(peer lifecycle chaincode queryinstalled -O json | jq -r '.installed_chaincodes | length')

if [ "$TEMP_COUNT" -eq 1 ]; 
then
    echo -e "\n\n* Copy package id"

    TEMP_VAR=$(peer lifecycle chaincode queryinstalled -O json | jq -r '.installed_chaincodes[0].package_id')

    export CC_PACKAGE_ID=$TEMP_VAR



    echo -e "\n\n* Approve"

    # peer lifecycle chaincode approveformyorg -n simplecc -v 1 -C  channel1 --sequence "$_SEQ" --init-required --package-id $CC_PACKAGE_ID --tls --cafile $ORDERER_CA
    peer lifecycle chaincode approveformyorg -n simplecc -v 1 -C  channel1 --sequence "$_SEQ" --package-id $CC_PACKAGE_ID --tls --cafile $ORDERER_CA



    echo -e "\n\n* Commit readiness"

    # peer lifecycle chaincode checkcommitreadiness -n simplecc -v 1 -C  channel1 --sequence "$_SEQ" --init-required
    peer lifecycle chaincode checkcommitreadiness -n simplecc -v 1 -C  channel1 --sequence "$_SEQ"



    echo -e "\n\n* Commit"

    # peer lifecycle chaincode commit -n simplecc -v 1 -C channel1 --sequence "$_SEQ" --init-required --tls --cafile $ORDERER_CA
    peer lifecycle chaincode commit -n simplecc -v 1 -C channel1 --sequence "$_SEQ" --tls --cafile $ORDERER_CA



    echo -e "\n\n* Check commited"

    peer lifecycle chaincode querycommitted -n simplecc  -C channel1


    # if [ "$_INIT" -eq "init"];
    # then
    #     echo -e "\n\n* Init"

    #     peer chaincode invoke --isInit  -n simplecc -C channel1 -c '{"Args":["init","a","100","b","200"]}' --tls --cafile $ORDERER_CA
    
    # fi

    echo -e "\n\n* Invoke to set super_admin"

    peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"InitLedger","Args":[]}' --tls --cafile $ORDERER_CA
    
    #1 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"CreateOrModify_Admin","Args":["123456", "THN1", "123456", "123456789013", "Admin1"]}' --tls --cafile $ORDERER_CA
    
    #2 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"Create_User","Args":["admin_THN1", "123456", "123456789014", "123456", "User1"]}' --tls --cafile $ORDERER_CA

    #3 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"Modify_User","Args":["admin_THN1", "123456", "123456789014", "User1"]}' --tls --cafile $ORDERER_CA
    
    #4 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"Create_Estate","Args":["admin_THN1", "123456", "200", "123456789014", "Badlapur", "100", "2021-12-15T20:34:33+05:30", "0"]}' --tls --cafile $ORDERER_CA
    
    #5 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"Modify_Estate","Args":["admin_THN1", "123456", "200", "", "101", "2021-12-15T20:34:33+05:30", "0"]}' --tls --cafile $ORDERER_CA
    
    #6 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"Verify_Estate","Args":["admin_THN1", "123456", "200", "1"]}' --tls --cafile $ORDERER_CA    

    #7 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"Verify_User","Args":["admin_THN1", "123456", "123456789014", "1"]}' --tls --cafile $ORDERER_CA
    
    #8 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"ChangeAvail_Estate","Args":["user_123456789014", "123456", "200", "true"]}' --tls --cafile $ORDERER_CA    

    #8 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"RequestToBuy_Estate","Args":["user_123456789014", "123456", "200", "99999", "2021-12-16T20:34:33+05:30"]}' --tls --cafile $ORDERER_CA

    #9 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"AcceptRequest_Estate","Args":["user_123456789014", "123456", "200", "123456789017", "2021-12-17T20:34:33+05:30"]}' --tls --cafile $ORDERER_CA

    #9 peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"ApproveSell_Estate","Args":["admin_THN1", "123456", "200", "2021-12-18T20:34:33+05:30"]}' --tls --cafile $ORDERER_CA

    echo -e "\n\n* Query to get super_admin"

    peer chaincode query -C channel1 -n simplecc  -c '{"Args":["GetValue","admin_super"]}'


    echo -e "\n\n* Query to get all data"

    peer chaincode query -C channel1 -n simplecc  -c '{"Args":["GetAll", "", ""]}'



    # echo -e "\n\n* Query to delete super_admin"

    # peer chaincode invoke -n simplecc -C channel1 -c '{"Function":"DeleteValue","Args":["admin_"]}' --tls --cafile $ORDERER_CA

elif [ "$TEMP_COUNT" -eq 0 ]
then
    echo -e "No chaincode intalled!"
else
    echo -e "Please clear previously installed packages and rerun."
fi

echo -e "\n\n* DONE"