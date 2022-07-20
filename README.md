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
- deploy a sample chaincode

**Query the network**
```
./main.sh chaincode query
```

This command will do a sample query on the network


# Todo
- [x] create a network: 2 orgs, 1 orderer
- [x] create channel script
- [x] join channel script
- [x] set anchor peer
- [x] create chaincode: insert, update, delete, query
- [ ] increase orgs: 4. update configuration
- [ ] use two chaincodes
- [ ] caliper
- [x] use single orgs endorsement policy
- [x] use couchdb & query with couchdb
