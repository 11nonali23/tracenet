'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
let ids = []

class CreateCampaignWorkload extends WorkloadModuleBase {

    constructor() {
        super();
        this.txIndex = 0;
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        ids = [];
    }

    async submitTransaction() {
        this.txIndex++;
        const randID = Math.floor(Math.random() * 1000)
        const assetID = randID.toString() + `_${this.workerIndex}_${this.txIndex}`;
        ids.push(assetID)

        console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateCampaign',
            invokerIdentity: 'peer0.obs0.tracenet.com',
            contractArguments: [assetID, 'Camp1', '"2022-05-02T15:02:40.628Z"', '"2023-05-02T15:02:40.628Z"'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
    async cleanupWorkloadModule() {
        for (let i = 0; i < ids.length; i++) {
            const assetID = ids[i];
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteCampaign',
                invokerIdentity: 'peer0.obs0.tracenet.com',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
            ids = []
        }
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */

function createWorkloadModule() {
    return new CreateCampaignWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;