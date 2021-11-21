# Start docker
```sudo systemctl start docker```

# Steps

- Installations
	- Git
	- Curl
	- Docker
	- Docker Compose
	- GO
	- JQ
	- Fabric Binaries
	
- Env Setup
	- add user to docker group
	- add go bin to PATH
	- copy fabric bin to user/local/bin
		- ```sudo cp ./fabric-samples/bin/* /usr/local/bin```

- cryptogen
	- ```cryptogen generate --config="./crypto-config.yaml" --output="crypto-material"```

- configtxgen
	- Create genesis block for "channel1"
		- ```configtxgen -profile OneOrgApplicationGenesis -outputBlock ./channel-artifacts/genesis.block -channelID channel1```
	<!-- - channel transaction
	- anchor peer transaction -->
- docker
	- write docker-compose.yaml
	- in root directory run
		- ```docker-compose -f ./docker-compose.yaml up -d```
		- check all containers
			- ```docker ps -a```
		- stop services
			- ```docker-compose -f ./docker-compose.yaml stop```
		- Stop and remove containers, networks, images, and volumes (for destroy everything)
			- ```docker-compose -f ./docker-compose.yaml down -v```
- connect to cli container
	- ```docker exec -it cli bash```

	- enroll (https://hyperledger-fabric.readthedocs.io/en/latest/create_channel/create_channel_participation.html)
	-	```
		export CHANNEL_NAME=channel1

		export OSN_TLS_CA_ROOT_CERT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

		export ADMIN_TLS_SIGN_CERT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

		export ADMIN_TLS_PRIVATE_KEY=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key
		```
	-	```osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/genesis.block -o "orderer.example.com:7053" --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY```

	- list channels
	
		```osnadmin channel list --channelID $CHANNEL_NAME -o "orderer.example.com:7053" --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY```

	- peer join channel (https://hyperledger-fabric.readthedocs.io/en/latest/commands/peerchannel.html)
	- 	```
		export CORE_PEER_TLS_ENABLED=true
		export CORE_PEER_LOCALMSPID="Org1MSP"
		export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
		export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
		export CORE_PEER_ADDRESS=peer1.org1.example.com:7051

		```
		```
		peer channel join -b ./channel-artifacts/genesis.block
		peer channel list
		```

	- 	```
		export CORE_PEER_TLS_ENABLED=true
		export CORE_PEER_LOCALMSPID="Org1MSP"
		export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/ca.crt
		export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
		export CORE_PEER_ADDRESS=peer2.org1.example.com:7061

		```
		```
		peer channel join -b ./channel-artifacts/genesis.block
		peer channel list
		```

	- anchor peer transaction, use configtxlator
		- env vars
			```
			export CORE_PEER_LOCALMSPID="Org1MSP"
			export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
			export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
			export CORE_PEER_ADDRESS=peer1.org1.example.com:7051
			export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
			```
		- The channel configuration block config_block.pb is stored in the channel-artifacts folder to keep the update process separate from other artifacts. (at this stage both genesis and latest channel config will be same because we haven't updated channel config)
			
			```
			# pb (extension for output): protobuf

			peer channel fetch config channel-artifacts/config_block.pb -o "orderer.example.com:7050" --ordererTLSHostnameOverride orderer.example.com -c channel1 --tls --cafile "$ORDERER_CA"
			
			```
		- decode the block from protobuf into a JSON object that can be read and edited
		
			```
			configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
			
			jq '.data.data[0].payload.data.config' config_block.json > config.json
			```
		- convert the channel configuration block into a streamlined JSON, config.json, that will serve as the baseline for our update. Because we don't want to edit this file directly, we will make a copy that we can edit. We will use the original channel config in a future step

			```cp config.json config_copy.json```
		- use the jq tool to add the Org1 anchor peer to the channel configuration
			
			```jq '.channel_group.groups.Application.groups.Org1MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer1.org1.example.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json```
		- convert both the original and modified channel configurations back into protobuf format and calculate the difference between them

			```
			configtxlator proto_encode --input config.json --type common.Config --output config.pb

			configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

			configtxlator compute_update --channel_id channel1 --original config.pb --updated modified_config.pb --output config_update.pb
			```
		- wrap the configuration update in a transaction envelope to create the channel configuration update transaction

			```
			configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json

			echo '{"payload":{"header":{"channel_header":{"channel_id":"channel1", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

			configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
			```
		- use the final artifact, config_update_in_envelope.pb, that can be used to update the channel
		- We can add the anchor peer by providing the new channel configuration to the peer channel update command. Because we are updating a section of the channel configuration that only affects Org1, other channel members do not need to approve the channel update.

			```
			export CORE_PEER_LOCALMSPID="Org1MSP"
			export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
			export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
			export CORE_PEER_ADDRESS=peer1.org1.example.com:7051
			export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
			```

			```peer channel update -f channel-artifacts/config_update_in_envelope.pb -c channel1 -o "orderer.example.com:7050"  --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA```
	
	- Chaincode
		- go to chaincode folder in cli file system
		- put chaincode file there
		- `go mod init github.com/chaincode` from chaincode directory itself
			- it will create go.mod
		- (https://golang.org/doc/tutorial/create-module)
		- `go mod tidy`
			- it will add and dependencies definde in chaincode file
		- `GO111MODULE=on go mod vendor`
			- it will install all dependencies and put them in vendor folder
		
			``````
		- package and install chaincode

			```
			export CORE_PEER_LOCALMSPID="Org1MSP"
			export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
			export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
			export CORE_PEER_ADDRESS=peer1.org1.example.com:7051

			export CORE_PEER_LOCALMSPID="Org1MSP"
			export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/ca.crt
			export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
			export CORE_PEER_ADDRESS=peer2.org1.example.com:7061
			```

			```
			# peer lifecycle chaincode package mycc.tar.gz -p . -n mycc --lang golang -v 1.0 -s -S
			peer lifecycle chaincode package ./packages/simple_cc_v1_v1.tar.gz -p . --label simplecc_v1_v1

			# peer lifecycle chaincode install mycc.tar.gz
			peer lifecycle chaincode install ./packages/simple_cc_v1_v1.tar.gz


			peer lifecycle chaincode queryinstalled
			```

		- approve chaincode and check commit readiness (admin of org)

			```
			export CORE_PEER_TLS_ENABLED=true
			export CORE_PEER_LOCALMSPID="Org1MSP"
			export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
			export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
			export CORE_PEER_ADDRESS=peer1.org1.example.com:7051

			export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

			export CC_PACKAGE_ID=simplecc_v1_v1:802ae430bc3e78eaf96bdac4c895cb33a57c54dacbeb7392f2e7a3071e77b858

			peer lifecycle chaincode approveformyorg -n simplecc -v 1 -C  channel1 --sequence 1 --init-required --package-id $CC_PACKAGE_ID --tls --cafile $ORDERER_CA

			peer lifecycle chaincode checkcommitreadiness -n simplecc -v 1 -C  channel1 --sequence 1 --init-required

			```
		- commit and check commited

			```
			peer lifecycle chaincode commit -n simplecc -v 1 -C channel1 --sequence 1 --init-required --tls --cafile $ORDERER_CA

			peer lifecycle chaincode querycommitted -n simplecc  -C channel1
			```
		- after commit
			1. Init the chaincode
				`peer chaincode invoke --isInit  -n simplecc -C channel1 -c '{"Args":["init","a","100","b","200"]}' --tls --cafile $ORDERER_CA`

			2. Query the chaincode
				`peer chaincode query -C channel1 -n simplecc  -c '{"Args":["query","a"]}'`

			3. Invoke the chaincode
				`peer chaincode invoke -C channel1 -n simplecc  -c '{"Args":["invoke","a","b","10"]}' --tls --cafile $ORDERER_CA`