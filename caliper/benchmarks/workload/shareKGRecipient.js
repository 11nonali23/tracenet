'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
let ids1 = []
let ids2 = []

class ShareKGVerificationWorkload extends WorkloadModuleBase {

    constructor() {
        super();
        this.txIndex = 0;
        this.campaignID = Math.floor(Math.random() * 1000).toString()
        this.KGId = Math.floor(Math.random() * 1000).toString()
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

        console.log(`Worker ${this.workerIndex}: Creating anonymized KG ${this.KGId}`);

        const shareKG = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'ShareAnonymizedKGForVerification',
            invokerIdentity: 'peer0.obs0.tracenet.com',
            contractArguments: [this.KGId, this.campaignID, "rec_id", "env", "sign"],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(shareKG);

        console.log(`Worker ${this.workerIndex}: Verifying the proof for the KG`);
        const verifyProof = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'VerifyProof',
            invokerIdentity: 'peer0.obs0.tracenet.com',
            contractArguments: [this.KGId, "equal-proof", "equal-proof"],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(verifyProof);
    }

    async submitTransaction() {
        this.txIndex++;
        const randID = Math.floor(Math.random() * 1000)
        const assetID1 = randID.toString() + `_${this.workerIndex}_${this.txIndex}`;
        const assetID2 = randID.toString() + `_${this.workerIndex}_${this.txIndex}`;

        console.log(`Worker ${this.workerIndex}: Share KG with recipient ${this.KGId}`);
        const shareKG = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CaliperShareAnonymizedKGWithRecipient',
            invokerIdentity: 'peer0.obs0.tracenet.com',
            contractArguments: [this.KGId, assetID1, assetID2, this.campaignID, "rec_id", "rec_env"],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(shareKG);
    }

    async cleanupWorkloadModule() {
        for (let i = 0; i < ids1.length; i++) {
            const assetID = ids1[i];
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteAnonymizedKG',
                invokerIdentity: 'peer0.obs0.tracenet.com',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
            ids1 = []
        }
        for (let i = 0; i < ids2.length; i++) {
            const assetID = ids2[i];
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteAnonymizedKG',
                invokerIdentity: 'peer0.obs0.tracenet.com',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
            ids2 = []
        }
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */

function createWorkloadModule() {
    return new ShareKGVerificationWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;