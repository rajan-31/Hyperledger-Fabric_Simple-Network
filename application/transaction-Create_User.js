'use strict';

const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');

async function main() {

    const wallet = await Wallets.newFileSystemWallet('./identity/user/User1@org1.example.com/wallet');

    const gateway = new Gateway();

    try {

        const userName = 'User1@org1.example.com';

        let connectionProfile = yaml.load(fs.readFileSync('./connection-org1.yaml', 'utf8'));

        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: { enabled:true, asLocalhost: true }
        };

        console.log('Connect to Fabric gateway.');

        await gateway.connect(connectionProfile, connectionOptions);

        console.log('Use network channel: channel1.');

        const network = await gateway.getNetwork('channel1');

        console.log('Use simplecc smart contract.');

        const contract = await network.getContract('simplecc');

        console.log('Submit Create_User transaction.');

        const issueResponse = await contract.submitTransaction('Create_User', "admin_THN1", "123456", "123456789017", "123456", "User2");

        console.log('Process Create_User transaction response: ', issueResponse);

        let json = JSON.parse(issueResponse.toString());

        console.log(`response: ${JSON.stringify(json, null, 2)}`);

        console.log('Transaction complete.');

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {

        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();

    }
}

main().then(() => {

    console.log('Create_User program complete.');

}).catch((e) => {

    console.log('Create_User program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});
