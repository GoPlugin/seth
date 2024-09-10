# Seth

Reliable and debug-friendly Ethereum client

[![Decoding tests](https://github.com/goplugin/seth/actions/workflows/test_decode.yml/badge.svg)](https://github.com/goplugin/seth/actions/workflows/test_decode.yml)
[![Tracing tests](https://github.com/goplugin/seth/actions/workflows/test_trace.yml/badge.svg)](https://github.com/goplugin/seth/actions/workflows/test_trace.yml)
[![API tests](https://github.com/goplugin/seth/actions/workflows/test_api.yml/badge.svg)](https://github.com/goplugin/seth/actions/workflows/test_api.yml)
[![CLI tests](https://github.com/goplugin/seth/actions/workflows/test_cli.yml/badge.svg)](https://github.com/goplugin/seth/actions/workflows/test_cli.yml)
[![Integration tests (testnets)](https://github.com/goplugin/seth/actions/workflows/test_decode_testnet.yml/badge.svg)](https://github.com/goplugin/seth/actions/workflows/test_decode_testnet.yml)
<br/>

*DISCLAIMER*: This repository is archived, there will be no new commits or support, all development is moved to https://github.com/goplugin/plugin-testing-framework

# Content
1. [Goals](#goals)
2. [Features](#features)
3. [Examples](#examples)
4. [Setup](#setup)
   1. [Building test contracts](#building-test-contracts)
   2. [Testing](#testing)
6. [Configuration](#config)
   1. [Simplified configuration](#simplified-configuration)
   2. [ClientBuilder](#clientbuilder)
   3. [Supported env vars](#supported-env-vars)
   4. [TOML configuration](#toml-configuration)
9. [Automated gas price estimation](#automatic-gas-estimator)
10. [DOT Graphs of transactions](#dot-graphs)
11. [Using multiple private keys](#using-multiple-keys)
12. [Experimental features](#experimental-features)
13. [Gas bumping for slow transactions](#gas-bumping-for-slow-transactions)
14. [CLI](#cli)
   1. [Manual gas price estimation](#manual-gas-price-estimation)
   2. [Block Stats](#block-stats)
   3. [Single transaction tracing](#single-transaction-tracing)
   4. [Bulk transaction tracing](#bulk-transaction-tracing)

## Goals

- Be a thin, debuggable and battle tested wrapper on top of `go-ethereum`
- Decode all transaction inputs/outputs/logs for all ABIs you are working with, automatically
- Simple synchronous API
- Do not handle `nonces` on the client side, trust the server
- Do not wrap `bind` generated contracts, small set of additional debug API
- Resilient: should execute transactions even if there is a gas spike or an RPC outage (failover)
- Well tested: should provide a suite of e2e tests that can be run on testnets to check integration

## Features

- [x] Decode named inputs
- [x] Decode named outputs
- [x] Decode anonymous outputs
- [x] Decode logs
- [x] Decode indexed logs
- [x] Decode old string reverts
- [x] Decode new typed reverts
- [x] EIP-1559 support
- [x] Multi-keys client support
- [x] CLI to manipulate test keys
- [x] Simple manual gas price estimation
- [ ] Fail over client logic
- [ ] Decode collided event hashes
- [x] Tracing support (4byte)
- [x] Tracing support (callTracer)
- [ ] Tracing support (prestate)
- [x] Tracing decoding
- [x] Tracing tests
- [ ] More tests for corner cases of decoding/tracing
- [x] Saving of deployed contracts mapping (`address -> ABI_name`) for live networks
- [x] Reading of deployed contracts mappings for live networks
- [x] Automatic gas estimator (experimental)
- [x] Block stats CLI
- [x] Check if address has a pending nonce (transaction) and panic if it does
- [x] DOT graph output for tracing
- [x] Gas bumping for slow transactions

You can read more about how ABI finding and contract map works [here](./docs/abi_finder_contract_map.md) and about contract store here [here](./docs/contract_store.md).

## Examples

Check [examples](./examples) folder

Lib provides a small amount of helpers for decoding handling that you can use with vanilla `go-ethereum` generated wrappers

```go
// Decode waits for transaction and decode all the data/errors
Decode(tx *types.Transaction, txErr error) (*DecodedTransaction, error)

// NewTXOpts returns a new sequential transaction options wrapper,
// sets opts.GasPrice and opts.GasLimit from seth.toml or override with options
NewTXOpts(o ...TransactOpt) *bind.TransactOpts

// NewCallOpts returns a new call options wrapper
NewCallOpts(o ...CallOpt) *bind.CallOpts
```

By default, we are using the `root` key `0`, but you can also use any of the private keys passed as part of `Network` configuration in `seth.toml` or ephemeral keys.

```go
// NewCallKeyOpts returns a new sequential call options wrapper from the key N
NewCallKeyOpts(keyNum int, o ...CallOpt) *bind.CallOpts

// NewTXKeyOpts returns a new transaction options wrapper called from the key N
NewTXKeyOpts(keyNum int, o ...TransactOpt) *bind.TransactOpts
```

Start `Geth` in a separate terminal, then run the examples

```sh
make GethSync
cd examples
go test -v
```

## Setup

We are using [nix](https://nixos.org/)

Enter the shell

```sh
nix develop
```

## Building test contracts

We have `go-ethereum` and [foundry](https://github.com/foundry-rs/foundry) tools inside `nix` shell

```sh
make build
```

## Testing

To run tests on a local network, first start it

```sh
make AnvilSync
```

Or use latest `Geth`

```sh
make GethSync
```

You can use default `hardhat` key `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80` to run tests

Run the [decode](./client_decode_test.go) tests

```sh
make network=Anvil root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test
make network=Geth root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test
```

Check other params in [seth.toml](./seth.toml), select any network and use your key for testnets

User facing API tests are [here](./client_api_test.go)

```sh
make network=Anvil root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test_api
make network=Geth root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test_api
```

CLI tests

```sh
make network=Anvil root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test_cli
make network=Geth root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test_cli
```

Tracing tests

```sh
make network=Anvil root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test_trace
make network=Geth root_private_key=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 test_trace
```

# Config

### Simplified configuration

If you do not want to set all the parameters, you can use a simplified progammatical configuration. Here's an example:

```go
cfg := seth.DefaultConfig("ws://localhost:8546", []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"})
client, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}
```

This config uses what we consider reasonable defaults, such as:
* 5 minute transaction confirmation timeout
* 1 minute RPC node dial timeout
* enabled EIP-1559 dynamic fees and automatic gas prices estimation (with 200 blocks history; will auto-disable itself if RPC doesn't support EIP-1559)
* tracing only of reverted transaction to console and DOT graphs
* checking of RPC node health on client creation
* no ephemeral keys

### ClientBuilder
You can also use a `ClientBuilder` to build a config programmatically. Here's an extensive example:

```go
client, err := builder.
    // network
    WithNetworkName("my network").
    WithRpcUrl("ws://localhost:8546").
    WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
    WithRpcDialTimeout(10*time.Second).
    WithTransactionTimeouts(1*time.Minute).
    // addresses
    WithEphemeralAddresses(10, 10).
    // tracing
    WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
    // protections
    WithProtections(true, true).
    // artifacts folder
    WithArtifactsFolder("some_folder").
    // nonce manager
    WithNonceManager(10, 3, 60, 5).
    // EIP-1559 and gas estimations
    WithEIP1559DynamicFees(true).
    WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
    WithGasPriceEstimations(false, 10, seth.Priority_Fast). 
	// gas bumping: retries, max gas price, bumping strategy function
    WithGasBumping(5, 100_000_000_000, PriorityBasedGasBumpingStrategyFn).	
    Build()

if err != nil {
    log.Fatal(err)
}
```

By default, it uses the same values as simplified configuration, but you can override them by calling the appropriate methods. Builder includes only options
that we thought to be most useful, it's not a 1:1 mapping of all fields in the `Config` struct. Therefore, if you need to set some more advanced options, you should create the `Config` struct directly,
use TOML config or manually set the fields on the `Config` struct returned by the builder.

### Supported env vars

Some crucial data is stored in env vars, create `.envrc` and use `source .envrc`, or use `direnv`

```sh
export SETH_LOG_LEVEL=info # global logger level
export SETH_CONFIG_PATH=seth.toml # path to the toml config
export SETH_NETWORK=Geth # selected network
export SETH_ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 # root private key

alias seth="SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go" # useful alias for CLI
```

Alternatively if you don't have a network defined in the TOML you can still use the CLI by providing these 2 key env vars:

```sh
export SETH_URL=https://rpc.fuji.testnet.anyswap.exchange
export SETH_CHAIN_ID=43113

go run cmd/seth/seth.go ... # your command
```

In that case you should still pass network name with `-n` flag.

### TOML configuration

Set up your ABI directory (relative to `seth.toml`)

```toml
abi_dir = "contracts/abi"
```

Setup your BIN directory (relative to `seth.toml`)

```toml
bin_dir = "contracts/bin"
```

Decide whether you want to generate any `ephemeral` keys:

```toml
# Set number of ephemeral keys to be generated (0 for no ephemeral keys). Each key will receive a proportion of native tokens from root private key's balance with the value equal to `(root_balance / ephemeral_keys_number) - transfer_fee * ephemeral_keys_number`.
ephemeral_addresses_number = 10
```

You can enable auto-tracing for all transactions meeting configured level, which means that every time you use `Decode()` we will decode the transaction and also trace all calls made within the transaction, together with all inputs, outputs, logs and events. Three tracing levels are available:

- `all` - trace all transactions
- `reverted` - trace only reverted transactions (that's default setting used if you don't set `tracing_level`)
- `none` - don't trace any transactions

Example:

```toml
tracing_level = "reverted"
```

Additionally, you can decide where tracing/decoding data goes to. There are three options:

- `console` - we will print all tracing data to the console
- `json` - we will save tracing data for each transaction to a JSON file
- `dot` - we will save tracing data for each transaction to a DOT file (graph)

```toml
trace_outputs = ["console", "json", "dot"]
```

For info on viewing DOT files please check the [DOT graphs](#dot-graphs) section below.

Example:
![image](./docs/tracing_example.png)
These two options should be used with care, when `tracing_level` is set to `all` as they might generate a lot of data.

If you want to check if the RPC is healthy on start, you can enable it with:

```toml
check_rpc_health_on_start = false
```

It will execute a simple check of transferring 10k wei from root key to root key and check if the transaction was successful.

You can add more networks like this:

```toml
[[Networks]]
name = "Fuji"
transaction_timeout = "30s"
# gas limit should be explicitly set only if you are connecting to a node that's incapable of estimating gas limit itself (should only happen for very old versions)
# gas_limit = 9_000_000
# hardcoded gas limit for sending funds that will be used if estimation of gas limit fails
transfer_gas_fee = 21_000
# legacy transactions
gas_price = 1_000_000_000
# EIP-1559 transactions
eip_1559_dynamic_fees = true
gas_fee_cap = 25_000_000_000
gas_tip_cap = 1_800_000_000
urls_secret = ["..."]
# if set to true we will dynamically estimate gas for every transaction (explained in more detail below)
gas_price_estimation_enabled = true
# how many last blocks to use, when estimating gas for a transaction
gas_price_estimation_blocks = 1000
# priority of the transaction, can be "fast", "standard" or "slow" (the higher the priority, the higher adjustment factor and buffer will be used for gas estimation) [default: "standard"]
gas_price_estimation_tx_priority = "slow"
```

If you don't we will use the default settings for `Default` network.

ChainID is not needed, as it's fetched from the node.

If you want to save addresses of deployed contracts, you can enable it with:

```toml
save_deployed_contracts_map = true
```

If you want to re-use previously deployed contracts you can indicate file name in `seth.toml`:

```toml
contract_map_file = "deployed_contracts_mumbai.toml"
```

Both features only work for live networks. Otherwise, they are ignored, and nothing is saved/read from for simulated networks.

### Automatic Gas Estimator

This section explains how to configure and understand the automatic gas estimator, which is crucial for executing transactions on Ethereum-based networks. Here’s what you need to know:

#### Configuration Requirements

Before using the automatic gas estimator, it's essential to set the default gas-related parameters for your network:

- **Non-EIP-1559 Networks**: Set the `gas_price` to define the cost per unit of gas if your network doesn't support EIP-1559.
- **EIP-1559 Networks**: If your network supports EIP-1559, set the following:
  - `eip_1559_dynamic_fees`: Enables dynamic fee structure.
  - `gas_fee_cap`: The maximum fee you're willing to pay per gas.
  - `gas_tip_cap`: An optional tip to prioritize your transaction within a block (although if it's set to `0` there's a high chance your transaction will take longer to execute as it will be less attractive to miners, so do set it).

These settings act as a fallback if the gas estimation fails. Additionally, always specify `transfer_gas_fee` for the fee associated with token transfers.

If you do not know if your network supports EIP-1559, but you want to give it a try it's recommended that you also set `gas_price` as a fallback. When we try to use EIP-1559 during gas price estimation, but it fails, we will fallback to using non-EIP-1559 logic. If that one fails as well, we will use hardcoded `gas_price` value.

#### How Gas Estimation Works

Gas estimation varies based on whether the network is a private Ethereum Network or a live network.

- **Private Ethereum Networks**: no estimation is needed. We always use hardcoded values.

For real networks, the estimation process differs for legacy transactions and those compliant with EIP-1559:

##### Legacy Transactions

1. **Initial Price**: Query the network node for the current suggested gas price.
2. **Priority Adjustment**: Modify the initial price based on `gas_price_estimation_tx_priority`. Higher priority increases the price to ensure faster inclusion in a block.
3. **Congestion Analysis**: Examine the last X blocks (as specified by `gas_price_estimation_blocks`) to determine network congestion, calculating the usage rate of gas in each block and giving recent blocks more weight.
4. **Buffering**: Add a buffer to the adjusted gas price to increase transaction reliability during high congestion.

##### EIP-1559 Transactions

1. **Tip Fee Query**: Ask the node for the current recommended tip fee.
2. **Fee History Analysis**: Gather the base fee and tip history from recent blocks to establish a fee baseline.
3. **Fee Selection**: Use the greater of the node's suggested tip or the historical average tip for upcoming calculations.
4. **Priority and Adjustment**: Increase the base and tip fees based on transaction priority (`gas_price_estimation_tx_priority`), which influences how much you are willing to spend to expedite your transaction.
5. **Final Fee Calculation**: Sum the base fee and adjusted tip to set the `gas_fee_cap`.
6. **Congestion Buffer**: Similar to legacy transactions, analyze congestion and apply a buffer to both the fee cap and the tip to secure transaction inclusion.

Understanding and setting these parameters correctly ensures that your transactions are processed efficiently and cost-effectively on the network.

Finally, `gas_price_estimation_tx_priority` is also used, when deciding, which percentile to use for base fee and tip for historical fee data. Here's how that looks:

```go
case Priority_Fast:
    baseFee = stats.GasPrice.Perc99
    historicalGasTipCap = stats.TipCap.Perc99
case Priority_Standard:
    baseFee = stats.GasPrice.Perc50
    historicalGasTipCap = stats.TipCap.Perc50
case Priority_Slow:
    baseFee = stats.GasPrice.Perc25
    historicalGasTipCap = stats.TipCap.Perc25
```

##### Adjustment factor

All values are multiplied by the adjustment factor, which is calculated based on `gas_price_estimation_tx_priority`:

```go
case Priority_Fast:
    return 1.2
case Priority_Standard:
    return 1.0
case Priority_Slow:
    return 0.8
```

For fast transactions we will increase gas price by 20%, for standard we will use the value as is and for slow we will decrease it by 20%.

##### Buffer percents

We further adjust the gas price by adding a buffer to it, based on congestion rate:

```go
case Congestion_Low:
    return 1.10, nil
case Congestion_Medium:
    return 1.20, nil
case Congestion_High:
    return 1.30, nil
case Congestion_VeryHigh:
    return 1.40, nil
```

For low congestion rate we will increase gas price by 10%, for medium by 20%, for high by 30% and for very high by 40%.

We cache block header data in an in-memory cache, so we don't have to fetch it every time we estimate gas. The cache has capacity equal to `gas_price_estimation_blocks` and every time we add a new element, we remove one that is least frequently used and oldest (with block number being a constant and chain always moving forward it makes no sense to keep old blocks).

It's important to know that in order to use congestion metrics we need to fetch at least 80% of the requested blocks. If that fails, we will skip this part of the estimation and only adjust the gas price based on priority.

For both transaction types if any of the steps fails, we fallback to hardcoded values.

### DOT graphs

There are multiple ways of visualising DOT graphs:

- `xdot` application [recommended]
- VSCode Extensions
- online viewers

### xdot

To install simply run `homebrew install xdot` and then run `xdot <path_to_dot_file>`. This tool seems to be the best for the job, since the viewer is interactive and supports tooltips, which in our case contain extra tracing information.

### VSCode Extensions

There are multiple extensions that can be used to view DOT files in VSCode. We recommend using [Graphviz Preview](https://marketplace.visualstudio.com/items?itemName=EFanZh.graphviz-preview). The downside is that it doesn't support tooltips.

### Goland

We were unable to find any (working) plugins for DOT graph visualization. If you do know any, please let us know.

### Online viewers

There's at least a dozen of them available, but none of them support tooltips and most can't handle our multi-line labels. These two are known to work, though:

- [Devtools/daily](https://www.devtoolsdaily.com/graphviz/)
- [Sketchviz](https://sketchviz.com/)

### Using multiple keys
If you want to use existing multiple keys (instead of ephemeral ones) you can pass them as part of the network configuration. In that case it's recommended to **not** read them from TOML file. If you need to read them for the filesystem/os it's best if you use environment variables.
Once you've read them in a safe manner you should programmatically add them to Seth's Config struct (which safe parts can be read from TOML file). You can either add them directly to `Network`, if it's already set up, or you can add them to `Networks` slice to the network you intend to use.

For example you could start by reading the TOML configuration first:
```go
cfg, err := seth.ReadCfg()
if err != nil {
    log.Fatal(err)
}
```

Then read the private keys in a safe manner. For example from a secure vault or environment variables:
```go
var privateKeys []string
var err error
privateKeys, err = some_utils.ReadPrivateKeysFromEnv()
if err != nil {
    log.Fatal(err)
}
```
and then add them to the `Network` you plan to use. Let's assume it's called `Sepolia`:
```go
for i, network := range cfg.Networks {
    if network.Name == "Sepolia" {
        cfg.Networks[i].PrivateKeys = privateKeys
    }
}
```

Or if you aren't using `[[Networks]]` in your TOML config and have just a single `Network`:
```go
cfg.Network.PrivateKeys = privateKeys
```

Or... you can use the convenience function `AppendPksToNetwork()` to have them added to both the `Network` and `Networks` slice:
```go
added := cfg.AppendPksToNetwork(privateKeys, "Sepolia")
if !added {
    log.Fatal("Network Sepolia not found in the config")
}
```

Finally, proceed to create a new Seth instance:
```go
seth, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}
```

A working example can be found [here](examples/example_test.go) as `TestSmokeExampleMultiKeyFromEnv` test.

Currently, there's no safe way to pass multiple keys to CLI. In that case TOML is the only way to go, but you should be mindful that if you commit the TOML file with keys in it, you should assume they are compromised and all funds on them are lost.

### Experimental features

In order to enable an experimental feature you need to pass its name in config. It's a global config, you cannot enable it per-network. Example:

```toml
# other settings before...
tracing_level = "reverted"
trace_outputs = ["console"]
experiments_enabled = ["slow_funds_return", "eip_1559_fee_equalizer"]
```

Here's what they do:

- `slow_funds_return` will work only in `core` and when enabled it changes tx priority to `slow` and increases transaction timeout to 30 minutes.
- `eip_1559_fee_equalizer` in case of EIP-1559 transactions if it detects that historical base fee and suggested/historical tip are more than 3 orders of magnitude apart, it will use the higher value for both (this helps in cases where base fee is almost 0 and transaction is never processed).

## Gas bumping for slow transactions
Seth has built-in gas bumping mechanism for slow transactions. If a transaction is not mined within a certain time frame (`Network`'s transaction timeout), Seth will automatically bump the gas price and resubmit the transaction. This feature is disabled by default and can be enabled by setting the `[gas_bumps] retries` to a non-zero number:
```toml
[gas_bumps]
retries = 5
```    

Once enabled, by default the amount, by which gas price is bumped depends on `gas_price_estimation_tx_priority` setting and is calculated as follows:
- `Priority_Fast`: 30% increase
- `Priority_Standard`: 15% increase
- `Priority_Slow`: 5% increase
- everything else: no increase

You can cap max gas price by settings (in wei):
```toml
[gas_bumps]
max_gas_price = 1000000000000
```

Once the gas price bump would go above the limit we stop bumping and use the last gas price that was below the limit.

How gas price is calculated depends on transaction type:
- for legacy transactions it's just the gas price
- for EIP-1559 transactions it's the sum of gas fee cap and tip cap
- for Blob transactions (EIP-4844) it's the sum of gas fee cap and tip cap and max fee per blob
- for AccessList transactions (EIP-2930) it's just the gas price

Please note that Blob and AccessList support remains experimental and is not tested.

If you want to use a custom bumping strategy, you can use a function with [GasBumpStrategyFn](retry.go) type. Here's an example of a custom strategy that bumps the gas price by 100% for every retry:
```go
var customGasBumpStrategyFn = func(gasPrice *big.Int) *big.Int {
    return new(big.Int).Mul(gasPrice, big.NewInt(2))
}
```

To use this strategy, you need to pass it to the `WithGasBumping` function in the `ClientBuilder`:
```go
var hundredGwei in64 = 100_000_000_000
client, err := builder.
    // other settings...
    WithGasBumping(5, hundredGwei, customGasBumpStrategyFn).
    Build()
```

Or set it directly on Seth's config:
```go
// assuming sethClient is already created
sethClient.Config.GasBumps.StrategyFn = customGasBumpStrategyFn
```

Since strategy function only accepts a single parameter, if you want to base its behaviour on anything else than that you will need to capture these values from the context, in which you define the strategy function. For example, you can use a closure to capture the initial gas price:
```go
gasOracleClient := NewGasOracleClient()

var oracleGasBumpStrategyFn = func(gasPrice *big.Int) *big.Int {
    // get the current gas price from the oracle
    suggestedGasPrice := gasOracleClient.GetCurrentGasPrice()

	// if oracle suggests a higher gas price, use it
    if suggestedGasPrice.Cmp(gasPrice) == 1 {
        return suggestedGasPrice
    }

	// otherwise bump by 100%
    return new(big.Int).Mul(gasPrice, big.NewInt(2))
}
```

Same strategy is applied to all types of transactions, regardless whether it's gas price, gas fee cap, gas tip cap or max blob fee.

When enabled, gas bumping is used in two places:
* during contract deployment via `DeployContract` function
* inside `Decode()` function

It is recommended to decrease transaction timeout when using gas bumping, as it will be effectively increased by the number of retries. So if you were running with 5 minutes timeout and 0 retries, you should set it to 1 minute and 5 retries
or 30 seconds and 10 retries.

Don't worry if while bumping logic executes previous transaction gets mined. In that case sending replacement transaction with higher gas will fail (because it is using the same nonce as original transaction) and we will retry waiting for the mining of the original transaction.

**Gas bumping is only applied for submitted transaction. If transaction was rejected by the node (e.g. because of too low base fee) we will not bump the gas price nor try to submit it, because original transaction submission happens outside of Seth.**

## CLI

You can either define the network you want to interact with in your TOML config and then refer it in the CLI command, or you can pass all network parameters via env vars. Most of the examples below show how to use the former approach.

### Manual gas price estimation

In order to adjust gas price for a transaction, you can use `seth gas` command

```sh
seth -n Fuji gas -b 10000 -tp 0.99
```

This will analyze last 10k blocks and give you 25/50/75/99th/Max percentiles for base fees and tip fees

`-tp 0.99` requests the 99th tip percentile across all the transaction in one block and calculates 25/50/75/99th/Max across all blocks

### Block Stats

If you need to get some insights into network stats and create a realistic load/chaos profile with simulators (`anvil` as an example), you can use `stats` CLI command

#### Define your network in `seth.toml`

Edit your `seth.toml`

```toml
[[networks]]
name = "MyCustomNetwork"
urls_secret = ["..."]

[block_stats]
rpc_requests_per_second_limit = 5
```

Then check the stats for the last N blocks

```sh
seth -n MyCustomNetwork stats -s -10
```

To check stats for the interval (A, B)

```sh
seth -n MyCustomNetwork stats -s A -e B
```

#### Pass all network parameters via env vars

If you don't have a network defined in the TOML you can still use the CLI by providing the RPC url via cmd arg.

Then check the stats for the last N blocks

```sh
seth -u "https://my-rpc.network.io" stats -s -10
```

To check stats for the interval (A, B)

```sh
seth -u "https://my-rpc.network.io" stats -s A -e B
```

Results can help you to understand if network is stable, what is avg block time, gas price, block utilization and transactions per second.

```toml
# Stats
perc_95_tps = 8.0
perc_95_block_duration = '3s'
perc_95_block_gas_used = 1305450
perc_95_block_gas_limit = 15000000
perc_95_block_base_fee = 25000000000
avg_tps = 2.433333333333333
avg_block_duration = '2s' 
avg_block_gas_used = 493233
avg_block_gas_limit = 15000000
avg_block_base_fee = 25000000000

# Recommended performance/chaos test parameters
duration = '2m0s'
block_gas_base_fee_initial_value = 25000000000
block_gas_base_fee_bump_percentage = '100.00% (no bump required)'
block_gas_usage_percentage = '3.28822000% gas used (no congestion)'
avg_tps = 3.0
max_tps = 8.0
```

### Single transaction tracing

You can trace a single transaction using `seth trace` command. Example with `seth` alias mentioned before:

```sh
seth -u "https://my-rpc.network.io" trace -t 0x4c21294bf4c0a19de16e0fca74e1ea1687ba96c3cab64f6fca5640fb7b84df65
```

or if you want to use a predefined-network:

```sh
seth -n=Geth trace -t 0x4c21294bf4c0a19de16e0fca74e1ea1687ba96c3cab64f6fca5640fb7b84df65
```

### Bulk transaction tracing

You can trace multiple transactions at once using `seth trace` command for a predefined network named `Geth`. Example:

```sh
seth -n=Geth trace -f reverted_transactions.json
```

or by passing all the RPC parameter with a flag:

```sh
seth -u "https://my-rpc.network.io" trace -f reverted_transactions.json
```

You need to pass a file with a list of transaction hashes to trace. The file should be a JSON array of transaction hashes, like this:

```json
[
  "0x...",
  "0x...",
  "0x...",
  ...
]
```

(Note that currently Seth automatically creates `reverted_transactions_<network>_<date>.json` with all reverted transactions, so you can use this file as input for the `trace` command.)
