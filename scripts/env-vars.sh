export FABRIC_LOGGING_SPEC=INFO
export FABRIC_CFG_PATH=./config

# Docker Compose
SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
export DOCKER_SOCK="${SOCK##unix://}"

# enroll orderer
export CHANNEL_NAME=channel1

# export OSN_TLS_CA_ROOT_CERT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# export ADMIN_TLS_SIGN_CERT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

# export ADMIN_TLS_PRIVATE_KEY=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

export ORDERER_CA=./crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=./crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=./crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key
