'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class VerifyProofWorkload extends WorkloadModuleBase {

    constructor() {
        super();
        this.campaignID = Math.floor(Math.random() * 1000).toString();
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        console.log(`Worker ${this.workerIndex}: Creating the campaign ${this.campaignID}`);

        const createCampaign = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateCampaign',
            invokerIdentity: 'peer0.obs0.tracenet.com',
            contractArguments: [this.campaignID, 'Camp1', '"2022-05-02T15:02:40.628Z"', '"2023-05-02T15:02:40.628Z"'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(createCampaign);

        for (let i = 0; i < 30; i++) {
            let KGID = Math.floor(Math.random() * 1000).toString()
            ids.push(KGID)
            console.log(`Worker ${this.workerIndex}: Creating a KG ${KGID}`);

            //should be an anonymized KG
            const shareKG = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'ShareKnowledgeGraph',
                invokerIdentity: 'peer0.obs0.tracenet.com',
                contractArguments: [KGID, this.campaignID, "abc", "10"],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(shareKG);
        }
    }

    async submitTransaction() {
        console.log(`Worker ${this.workerIndex}: Verifying the proof`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'VerifyProof',
            invokerIdentity: 'peer0.obs0.tracenet.com',
            contractArguments: [],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
    async cleanupWorkloadModule() {
        //nothing here
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */

function createWorkloadModule() {
    return new VerifyProofWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;