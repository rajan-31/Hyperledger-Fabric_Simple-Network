'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');

// Main program function
async function main() {

    // A wallet stores a collection of identities for use
    const wallet = await Wallets.newFileSystemWallet('./identity/user/User1@org1.example.com/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    // Main try/catch block
    try {

        // Specify userName for network access
        // const userName = 'isabella.issuer@magnetocorp.com';
        const userName = 'User1@org1.example.com';

        // Load connection profile; will be used to locate a gateway
        let connectionProfile = yaml.load(fs.readFileSync('./connection-org1.yaml', 'utf8'));

        // Set connection options; identity and wallet
        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: { enabled:true, asLocalhost: true }
        };

        // Connect to gateway using application specified parameters
        console.log('Connect to Fabric gateway.');

        await gateway.connect(connectionProfile, connectionOptions);

        // Access PaperNet network
        console.log('Use network channel: channel1.');

        const network = await gateway.getNetwork('channel1');

        // Get addressability to commercial paper contract
        console.log('Use simplecc smart contract.');

        const contract = await network.getContract('simplecc');

        // issue commercial paper
        console.log('Submit CreateOrModify_Admin transaction.');

        const issueResponse = await contract.submitTransaction('CreateOrModify_Admin', "123456", "THN2", "123456", "123456789015");

        // process response
        console.log('Process CreateOrModify_Admin transaction response: ', issueResponse);

        // let json = JSON.parse(issueResponse.toString());

        // console.log(`response: ${json}`);        

        console.log('Transaction complete.');

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {

        // Disconnect from the gateway
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();

    }
}

main().then(() => {

    console.log('CreateOrModify_Admin program complete.');

}).catch((e) => {

    console.log('CreateOrModify_Admin program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});
