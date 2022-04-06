const express = require("express");
const router = express.Router();

router.post("/api/Create_User", async (req, res) => {

    const _data = req.body;
    const _username = _data.username,
        _password = _data.password,
        uid = _data.uid,
        newPass = _data.new_pass,
        newName = _data.new_name;
    
    try {
        // console.log('Submit Create_User transaction.');
        const issueResponse = await contract.submitTransaction('Create_User', _username, _password, uid, newPass, newName);

        // console.log('Process Create_User transaction response: ', issueResponse);
        let json = JSON.parse(issueResponse.toString());

        // console.log(`response: ${JSON.stringify(json, null, 2)}`);
        // console.log('Transaction complete.');

        res.status(200).json(json);
    } catch (error) {
        console.log(error);
        res.status(551).send(error.toString());
    }
});

// router.post("/api/GetValue", async (req, res) => {
//     const key = req.body.key;

//     try {
//         const issueResponse = await contract.submitTransaction('GetValue', key);
//         let json = JSON.parse(issueResponse.toString());

//         res.status(200).json(json);
//     } catch (error) {
//         console.log(error);
//         res.status(551).send(error.toString());
//     }
// });

router.post("/api/GetValue", (req, res) => {
    const key = req.body.key;

    contract.submitTransaction('GetValue', key).then((payload)=>{
        let json = JSON.parse(payload.toString());
        console.log(json.password);
        res.status(200).json(json);
    }).catch((error) => {
        console.log(error);
        res.status(551).send(error.toString());
    });
});

module.exports = router;