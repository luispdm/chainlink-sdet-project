# Chainlink SDET Project

## Requirements
To run the tests, you must have [go](https://golang.org/) or [Docker](https://www.docker.com/) installed.
You also need a valid WebSocket address to make requests to the blockchain. You can leverage services like [Alchemy](https://www.alchemy.com/) or [Infura](https://infura.io/).

## Configuration
As of 06/2020, the WebSocket address and the test parallelization are configurable.
You must either create a file named `config.yml` under the `config` directory or pass-in the environment variables.

Config file example:
```yaml
wss: "myWebSocketAddr"
parallel: true
```
Environment variables example:
```bash
PARALLEL=false WSS=myAddr go test
```
**The environment variables take precedence over the config file.**
**Without a WebSocket address, the tests will fail.**

## Instructions
To run the tests, you can use the provided `Makefile` recipes, use your own `go` executable or build and run the Docker container.

### Local
Examples:
```bash
make go-test                          # Runs "go test -v"
make compile                          # Compiles the test binary
make binary                           # Runs the compiled binary (equivalent to "make go-test")
WSS=myAddr PARALLEL=false make binary # Same as "make binary" but with env vars
```

### Docker
The process of building the image is a bit slow because Docker downloads the dependencies and compiles the test binary. After that, each execution is faster because the container runs the compiled binary.

Examples:
```bash
make build             # Builds the Docker image
make run               # Runs the container with the compiled test binary
PARALLEL=true make run # Same as "make run" but with env vars
```

The requirement of the WebSocket is valid for the container as well.

## Choices
The Solidity source code at `contracts/ethereum/v0.6.6/src/Contract.sol` comes from [this BTC/USD feed contract](https://etherscan.io/address/0xf570deefff684d964dc3e15e1f9414283e3f7419#code).
The abi and bin files stored at `contracts/ethereum/v0.6.6` were generated with [solc 0.6.6](https://github.com/ethereum/solidity/releases/tag/v0.6.6). The go objects stored at `contracts/ethereum` were generated with [abigen](https://geth.ethereum.org/docs/dapp/native-bindings).

The tests leverage [testify](https://github.com/stretchr/testify) for assertions. The same library is also used for `Setup` and `TearDown` of suites and tests.

The tests are table-driven with the following variables:
1. test name
2. contract address
3. deviation threshold
4. number of rounds

I think having the number of rounds and threshold configurable per-test rather than global makes the tests more customizable.

There's a maximum number of rounds allowed stored in the constant `MAX_ROUNDS` at `main_test.go`.

## Observations
- The maximum number of rounds could have been parameterized via config file or env var. The creator of this repo forgot to do so 😅.
- There's no usage of goroutines or thread-safe constructs like `sync.Map` because, after a couple of benchmarks, it's been observed that the optimization brought by the concurrency was poor (and the overhead increased). I preferred the simplicity over the speed.
- I have noticed that the BTC/USD feed and the LINK/USD feed used in this test have 15 nodes. However, while testing, only the data of 14 nodes is returned. I do not understand why.