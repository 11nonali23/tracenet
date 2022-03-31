#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh
. $SCRIPTS_DIR/utils/environment.sh




MODE=$1
CHAINCODE_NAME=$2
CHANNEL_NAME=$3

if [ "$MODE" == "package" ]; then
  packageChaincode $CHAINCODE_NAME
elif [ "$MODE" == "install" ]; then
  installChaincode $CHAINCODE_NAME $CHANNEL_NAME "rec" 0 0
  installChaincode $CHAINCODE_NAME $CHANNEL_NAME "obs" 0 0
elif [ "$MODE" == "approve" ]; then
  approveForMyOrg $CHAINCODE_NAME $CHANNEL_NAME "rec" 0 0
  approveForMyOrg $CHAINCODE_NAME $CHANNEL_NAME "obs" 0 0
elif [ "$MODE" == "commit" ]; then
  commitChaincode $CHAINCODE_NAME $CHANNEL_NAME "rec,obs" 1 1
fi