# DERO Tools

## Cold Wallet
Allows generating in a simple way and without connection, a `DERO address`, the associated `seed`, and the valid TX Registration.
To validate the account creation on blockchain, you must propagate the `TX Registration hex` using [API](https://docs.dero.io/rtd_pages/dev_rpcapistargate.html#send-raw-transaction).


### Changes

Faster registration process. Display ready-to-use `curl` command. Only the daemon address needs to be replaced.

### Building

```
go mod tidy
go build -o coldwallet ./cmd/coldwallet
```

### Example
Seed language 0 (english)

`./coldwallet --language=0`

Example output
```
Language seed selected: English
Version: 3.5.3-140.DEROHE.STARGATE+13062023
Generating new random wallet...
Address: dero1address
Seed   : seed phrase appears here
TX data: 0100000123456789ABCDEF

curl -X POST <node_address>:10102/json_rpc -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":1,"method":"sendrawtransaction","params":{"tx_as_hex":"0100000123456789ABCDEF"}}'
```