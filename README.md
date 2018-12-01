Inspired by [evm-tools]. (https://github.com/CoinCulture/evm-tools)(But this project has stopped development and is now unable to compile).

This is a guide to understanding the EVM, its relationship with solidity, and how to use some debugging tools.

# Requirements
- evmfun has been tested and is known to run on Linux/Ubuntu, macOS and Windows(10). It will likely work fine on most OS.
- Go 1.11 or newer.

# Installation

```sh
go get github.com/mohanson/evmfun/cmd/evmfun
```

# Feat: disassemble bytecode

Here is some very simple bytecode I wrote:

```
0x6005600401
```

To disassemble, run `evm disasm 0x6005600401`, which produces:

```
[0000] PUSH1 0x05
[0002] PUSH1 0x04
[0004] ADD
```

# Feat: execute bytecode

I found a code in the test library: [https://github.com/ethereum/tests/blob/develop/VMTests/vmSystemOperations/return0.json](https://github.com/ethereum/tests/blob/develop/VMTests/vmSystemOperations/return0.json). Execute the test now

```
evm exec 0x603760005360005160005560016000f3
```

```
[Many Outputs]
...

Return = 0x37
```

Perfect!
