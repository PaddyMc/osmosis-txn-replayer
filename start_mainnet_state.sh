#!/bin/bash
set -e

OSMOSIS_HOME=$HOME/.osmosisd

SNAPSHOT_URL=$(curl -sL https://snapshots.osmosis.zone/latest)
RPC_URL=https://rpc.osmosis.zone
ADDRBOOK_URL="https://rpc.osmosis.zone/addrbook"
GENESIS_URL=https://github.com/osmosis-labs/osmosis/raw/main/networks/osmosis-1/genesis.json

rm -rf $HOME/.osmosisd/

osmosisd init test --chain-id testing --home=$HOME/.osmosisd

# Copy genesis
echo -e "\nDownloading genesis file..."
wget $GENESIS_URL -O $OSMOSIS_HOME/config/genesis.json
echo ✅ Genesis file downloaded successfully.

# Download latest snapshot
echo -e "\nDownloading latest snapshot..."
wget -O - $SNAPSHOT_URL | lz4 -d | tar -C $OSMOSIS_HOME/ -xf -
echo -e ✅ Snapshot downloaded successfully.

# Run the node
osmosisd start --home=$HOME/.osmosisd
