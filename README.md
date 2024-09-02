### osmosis-txn-replayer

Omsosis txn replayer is a testing tool to replay transactions from mainnet to a local edgenet to simulate mainnet after an upgrade.

### How to use the script

- Run `./start_mainnet_state.sh`
- Before the chain gets fully synced stop the node
- Run this script e.g `go run main.go`
- On the current version e.g v25 run `osmosisd in-place-testnet osmosis-1 osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj --trigger-testnet-upgrade v26`
- Then when the node stope with `3:14PM ERR CONSENSUS FAILURE!!! err="UPGRADE \"v26\" NEEDED at height: 20246428: " `
- Stop the node and recompile the lastest version e.g v26
- Run this script e.g `go run main.go` again
- Then run  `osmosisd start --home=$HOME/.osmosisd --rpc.unsafe --grpc.enable --grpc-web.enable --api.enabled-unsafe-cors --api.enable`

### Beware of dragons
- Ensure the edgenet is running slower than mainnet by setting `timeout_commit = 3s` in the config.toml
- Some txns will fail most should pass
