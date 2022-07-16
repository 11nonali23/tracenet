#!/bin/bash

. $PWD/settings.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=INFO
export CHAINCODE_NAME="main"

function initialize() {
    $SCRIPTS_DIR/init.sh "orgs"
    sleep 1
    $SCRIPTS_DIR/init.sh "system-genesis-block"
}

function networkUp() {
    $SCRIPTS_DIR/network.sh "start" $LOG_LEVEL
}

function networkDown() {
    $SCRIPTS_DIR/network.sh "stop" $LOG_LEVEL
}

function clear() {
    $SCRIPTS_DIR/network.sh "clear"
}

function createChannel() {
    $SCRIPTS_DIR/channel.sh "create-tx" $CHANNEL_NAME
    sleep 3
    $SCRIPTS_DIR/channel.sh "create" $CHANNEL_NAME
}

function joinChannel() {
    $SCRIPTS_DIR/channel.sh "join" $CHANNEL_NAME
}

function packageChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "package" $CHAINCODE_NAME
}

function installChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "install" $CHAINCODE_NAME $CHANNEL_NAME
    $SCRIPTS_DIR/deployChaincode.sh "install" $CHAINCODE_NAME $CHANNEL_NAME
}

function approveChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "approve" $CHAINCODE_NAME $CHANNEL_NAME
    $SCRIPTS_DIR/deployChaincode.sh "approve" $CHAINCODE_NAME $CHANNEL_NAME
}

function commitChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "commit" $CHAINCODE_NAME $CHANNEL_NAME
}

# function checkCommitted() {
#     $SCRIPTS_DIR/verifyChaincode.sh "committed" $CHANNEL_NAME "rec" 0 0
#     $SCRIPTS_DIR/verifyChaincode.sh "committed" $CHANNEL_NAME "obs" 0 0
# }

# function checkInstalled() {
#     $SCRIPTS_DIR/verifyChaincode.sh "installed" "rec" 0 0
#     $SCRIPTS_DIR/verifyChaincode.sh "installed" "obs" 0 0
# }

# function checkCommitReadliness() {
#     $SCRIPTS_DIR/verifyChaincode.sh "ready" $CHAINCODE_NAME $CHANNEL_NAME "rec" 0 0
#     $SCRIPTS_DIR/verifyChaincode.sh "ready" $CHAINCODE_NAME $CHANNEL_NAME "obs" 0 0
# }

# function listChaincode() {
#     $SCRIPTS_DIR/verifyChaincode.sh "list" $CHANNEL_NAME "rec" 0 0
#     $SCRIPTS_DIR/verifyChaincode.sh "list" $CHANNEL_NAME "obs" 0 0
# }

function invokeChaincodeInit() {
    id="c${RANDOM}"
    startDate="2022-05-02T15:02:40.628Z"
    endDate="2023-05-02T15:02:40.628Z"
    fcnCall='{"function":"'CreateCampaign'","Args":["'${id}'","'Campaign1'","'Rec0'","'${startDate}'","'${endDate}'"]}'
    $SCRIPTS_DIR/chaincodeOperation.sh $CHAINCODE_NAME $CHANNEL_NAME "rec,obs" 1 1 $fcnCall
}

function queryChaincode() {
    fcnCall='{"function":"'ReadAllCampaigns'","Args":[]}'

    $SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "rec" 1 1 $fcnCall
}

function initCaliper() {
    $SCRIPTS_DIR/caliper.sh "init" $CALIPER_VERSION $FABRIC_VERSION
}

function caliperLaunchManager() {
    $SCRIPTS_DIR/caliper.sh "launch-manager" $CALIPER_VERSION $FABRIC_VERSION $CALIPER_WORKSPACE $CALIPER_NETWORK_CONFIG $CALIPER_BENCH_CONFIG
}

function clearCaliper() {
    $SCRIPTS_DIR/caliper.sh "clear"
}

MODE=$1

if [ $MODE = "network" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "up" ]; then
        initialize
        networkUp
    elif [ $SUB_MODE = "down" ]; then
        networkDown
        clear
    elif [ $SUB_MODE = "restart" ]; then
        networkDown
        clear
        initialize
        networkUp
        createChannel
        joinChannel
        packageChaincode
        installChaincode
        approveChaincode
        commitChaincode
        invokeChaincodeInit
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "channel" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "create" ]; then
        createChannel
    elif [ $SUB_MODE = "join" ]; then
        joinChannel
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "chaincode" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "package" ]; then
        packageChaincode
    elif [ $SUB_MODE = "install" ]; then
        installChaincode
    elif [ $SUB_MODE = "approve" ]; then
        approveChaincode
    elif [ $SUB_MODE = "commit" ]; then
        commitChaincode

    # elif [ $SUB_MODE = "check" ]; then
    #     SUB_SUB_MODE=$3

    #     if [ $SUB_SUB_MODE = "installed" ]; then
    #         checkInstalled
    #     elif [ $SUB_SUB_MODE = "ready" ]; then
    #         checkCommitReadliness
    #     elif [ $SUB_SUB_MODE = "committed" ]; then
    #         checkCommitted
    #     elif [ $SUB_MODE = "list" ]; then
    #         listChaincode
    #     else
    #         echo "Unsuported '$MODE $SUB_MODE $SUB_SUB_MODE' command."
    #     fi

    elif [ $SUB_MODE = "invoke-init" ]; then
        invokeChaincodeInit
    elif [ $SUB_MODE = "query" ]; then
        queryChaincode
    elif [ $SUB_MODE = "reinstall" ]; then
        packageChaincode
        installChaincode
        approveChaincode
        commitChaincode
    else
        echo "Unsupported '$MODE $SUB_MODE' command."
    fi

elif [ $MODE = "caliper" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "init" ]; then
        initCaliper
    elif [ $SUB_MODE = "launch-manager" ]; then
        caliperLaunchManager
    elif [ $SUB_MODE = "clear" ]; then
        clearCaliper
    else
        echo "Unsupported '$MODE $SUB_MODE' command."
    fi
else
    echo "Unsupported $MODE command."
fi