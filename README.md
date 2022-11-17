# Tracenet
An Hyperledger Fabric network. It has 2 orgs and 1 orderer.
For the certificates generation I used cryptogen, since it's a testing network.
Each peer has a couchdb instance.

# Prerequisites 
- [go](https://go.dev) --> It's the only language supported
- [docker](https://www.docker.com)
- [docker-compose](https://docs.docker.com/compose/)
- [Hyperledger Fabric 2.2](https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html)
- [npm](https://www.npmjs.com) --> only for testing with caliper
- [jq](https://stedolan.github.io/jq/)

# Running and working on the netowrk

**Script permissions**
```
sudo chmod 755 main.sh
sudo chmod 755 settings.sh
sudo chmod -R 755 scripts/
```
**Run the network**
```
./main.sh network restart
```
This command will:
- down the network and clean
- initialization with cryptogen
- start the network with docker compose
- create a channel named "mychannel"
- join the channel by the two orgs (Org1, Org2)
- set the anchor peer
- deploy the chaincode

# Testing the network

The framework used for the tests is Caliper. Test are run after initializing the ledger with 2000 transactions. Then every workload has 2000 tx with 200tps.
```
./main.sh caliper init
./main.sh caliper launch
```
