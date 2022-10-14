'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class VerifyProofWorkload extends WorkloadModuleBase {

    constructor() {
        super();
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
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