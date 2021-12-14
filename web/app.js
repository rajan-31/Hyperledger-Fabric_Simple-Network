const express = require("express")
const app = express();

const path = require("path")

const indexRoutes = require("./routes/index"),
    apiRoutes = require("./routes/api");

/* static content */
app.use("/public", express.static(path.join(__dirname, "public"), {
    etag: false,
    // maxAge: 1000 * 60 * 60   // 1 hr
}));

require('dotenv').config();  // loading environment variables

app.use(express.urlencoded({extended: true}));
app.use(express.json());

app.use(indexRoutes);
app.use(apiRoutes);

//=====================================

contract = null;

const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');

const connectionSetup = async ()=> {
    const wallet = await Wallets.newFileSystemWallet(path.join(__dirname, "..", "application", "identity/user/User1@org1.example.com/wallet"));
    const gateway = new Gateway();

    try {
        const userName = 'User1@org1.example.com';
        let connectionProfile = yaml.load(fs.readFileSync(path.join(__dirname, "..", "application", "connection-org1.yaml"), 'utf8'));
        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: { enabled:true, asLocalhost: true }
        };

        console.log('Connect to Fabric gateway.');
        await gateway.connect(connectionProfile, connectionOptions);

        const network = await gateway.getNetwork('channel1');
        contract = await network.getContract('simplecc');

        console.log("Connection is set :)");
    } catch (error) {

        console.log(`Error creating Gateway: ${error}`);
        console.log(error.stack);

    }
}

connectionSetup();

//=====================================

app.get('/*', function(req, res){
    res.status(404).send(`
    <div style="text-align:center;">
        <h1>Error 404</h1>
        <h3>Are you lost?!</h3>
        <a href="/">Go Home</a>
    </div>
    `)
});

const port = process.env.PORT

app.listen(port, function(){
    console.log("Server is running on port => " + port);
});