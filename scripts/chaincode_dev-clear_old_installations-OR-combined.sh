#!/bin/bash

# run in user shell

echo "* Kill running chaincode container"

TEMP_CONTAINER=$(docker ps -q --filter ancestor=$(docker images --filter=reference='dev-peer1*' -q))

if [ "$TEMP_IMAGE" != "" ];
then
    docker kill $(docker ps -q --filter ancestor=$(docker images --filter=reference='dev-peer1*' -q))

else
    echo "No Container found."
fi



echo "\n\n* Remove installed chaincodes from peer1"

docker exec peer1.org1.example.com sh -c 'rm -rf /var/hyperledger/production/lifecycle/chaincodes/*.tar.gz'

echo "\n\n* Remove old packages from cli"

docker exec cli sh -c 'rm -rf /opt/gopath/src/github.com/chaincode/packages/*.tar.gz'



echo "\n\n* Remove chaincode image"

TEMP_IMAGE=$(docker images --filter=reference='dev-peer1*' -q)

if [ "$TEMP_IMAGE" != "" ];
then
    docker rmi -f $(docker images --filter=reference='dev-peer1*' -q)

else
    echo "No image found."
fi



echo "\n\n* Restart using compose"

docker-compose restart



echo "\n\n* DONE"



echo "\n\n--------------------------------------------------------"

if [ "$1" = "combined" ];
then
    docker exec cli sh -c 'sh ./scripts/chaincode_dev-main.sh . $(($(peer lifecycle chaincode querycommitted -n simplecc  -C channel1 -O json | jq -r '.sequence') + 1))'

fi