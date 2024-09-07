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

# Wait for at least 1 block
#  killall osmosisd

# Make sure you're on v25
# /home/cryptomatictrader/github.com/osmosis-labs/osmosis.v25/build/osmosisd in-place-testnet edgenet osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj --trigger-testnet-upgrade v26

# Run to run migrations
# /home/cryptomatictrader/github.com/osmosis-labs/osmosis.cleanV26/build/osmosisd start --home=$HOME/.osmosisd --p2p.persistent_peers "" --p2p.seeds "" --rpc.unsafe --grpc.enable --grpc-web.enable

