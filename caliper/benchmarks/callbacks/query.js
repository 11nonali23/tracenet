'use strict';

module.exports.info = 'Sample Workload';

const { v1: uuidv4 } = require('uuid');

let txCountPerBatch;
let blockchain, context;

module.exports.init = function (blockchain, context, args) {

    if (!args.hasOwnProperty('txCountPerBatch')) {
        args.txCountPerBatch = 1;
    }
    txCountPerBatch = args.txCountPerBatch;
    blockchain = blockchain;
    context = context;

    return Promise.resolve();
};

function generateWorkload() {
    let workload = [];
    for (let i = 0; i < txCountPerBatch; i++) {
        workload.push({
            chaincodeFunction: 'CreateAsset',
            chaincodeArguments: [uuidv4(), "white", 100, "Andrea", 90],
        });
    }
    return workload;
}

module.exports.run = function () {
    try {
        let args = generateWorkload();
        return blockchain.invokeSmartContract(context, 'sample', '1.0', args);
    } catch (err) {
        console.log(err)
    }

};

module.exports.end = function () {
    return Promise.resolve();
};