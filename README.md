### osmosis-txn-replayer

Omsosis txn replayer is a testing tool to replay transactions from mainnet to a local edgenet to simulate mainnet after an upgrade.

### How to use the script

1. Run `./start_mainnet_state.sh`
2. Before the chain gets fully synced stop the node
3. Run this script e.g `go run main.go`. Leave it running until after step 5 below.
4. On the current version e.g v25 run `osmosisd in-place-testnet osmosis-1 osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj --trigger-testnet-upgrade v26`. Note the in-place-testnet is created with chain id `osmosis-1`.
5. Then when the node stope with `3:14PM ERR CONSENSUS FAILURE!!! err="UPGRADE \"v26\" NEEDED at height: 20246428: " `
6. Stop the node and recompile the lastest version e.g v26
7. Run this script e.g `go run main.go`. Again, leave it running while the node is running in step 8.
8. Then run  `osmosisd start --home=$HOME/.osmosisd --rpc.unsafe --grpc.enable --grpc-web.enable --api.enabled-unsafe-cors --api.enable`

### Logs 
The logs output block height, type of message and txn hash if sucessful, or an error if theres a failure.
```
2024/09/03 08:54:50 main.go:130: Processed block at height: 20262208
2024/09/03 08:54:53 main.go:144: Message 0: Type: *types.MsgSend
2024/09/03 08:54:53 main.go:157: Transaction replayed successfully. Hash: E85E55E35225C5B7DBE14715CA8F99AC53361477D46A33CBCAE4363E8A27E06D
2024/09/03 08:54:53 main.go:144: Message 0: Type: *types.MsgWithdrawDelegatorReward
2024/09/03 08:54:53 main.go:157: Transaction replayed successfully. Hash: 6EA638CA805C14B438FB52923CFE713186A323060A400AD2C237075ECDA0799D
2024/09/03 08:54:53 main.go:144: Message 0: Type: *types.MsgExecuteContract
2024/09/03 08:54:53 main.go:157: Transaction replayed successfully. Hash: 836A7FDC985802E9EC7F8C6FECFAC90C9BE1395D3C33D99E5095D7FC65F16BAE
2024/09/03 08:54:53 main.go:144: Message 0: Type: *types.MsgExecuteContract
2024/09/03 08:54:53 main.go:144: Message 1: Type: *types.MsgExecuteContract
2024/09/03 08:54:53 main.go:157: Transaction replayed successfully. Hash: 7D10867674C5E315CFEB117B18C6F215710AB4A83961B447FB6C2D3819441E82
2024/09/03 08:54:53 main.go:144: Message 0: Type: *types.MsgExecuteContract
2024/09/03 08:54:53 main.go:157: Transaction replayed successfully. Hash: 839FF6966F927E8AD4B8CFD0ED99892C5E435E895BBFA2F00AD4947D36B1010E
2024/09/03 08:54:53 main.go:130: Processed block at height: 20262209
2024/09/03 08:54:56 main.go:144: Message 0: Type: *types.MsgExecuteContract
2024/09/03 08:54:56 main.go:157: Transaction replayed successfully. Hash: C7B252AF9BC8167A14112FCB8B4659390F7132D8CDD2F0ED799C3673DE80939B
```

### Beware of dragons
- Ensure the edgenet is running slower than mainnet by setting `timeout_commit = 3s` in the config.toml
- Some txns will fail most should pass. You can double-check the transaction status by `osmosisd q tx [hash]`
