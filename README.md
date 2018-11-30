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

To disassemble, run `echo 0x6005600401 | evmfun disasm`, which produces:

```
0      PUSH1   => 05
2      PUSH1   => 04
4      ADD
```
