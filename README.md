Inspired by [evm-tools](https://github.com/CoinCulture/evm-tools). (But this project has stopped development and is now unable to compile).

This is a guide to understanding the EVM, its relationship with solidity, and how to use some debugging tools.

# Requirements
- evm has been tested and is known to run on Linux/Ubuntu, macOS and Windows(10). It will likely work fine on most OS.
- [Go](https://golang.org/dl/) 1.11 or newer.

# Installation

```sh
go get github.com/mohanson/evm/cmd/evm
```

# Feat: disassemble bytecode

Here is some very simple bytecode I wrote:

```
0x6005600401
```

To disassemble, run `evm disasm -code 0x6005600401`, which produces:

```
[0000] PUSH1 0x05
[0002] PUSH1 0x04
[0004] ADD
```

# Feat: execute bytecode

Let's use a test case to demonstrate: [return0.json](https://github.com/ethereum/tests/blob/develop/VMTests/vmSystemOperations/return0.json). To execute the test, run

```sh
$ evm exec -code 0x603760005360005160005560016000f3
```

```
[Many Outputs]
...

Return = 0x37
```

Perfect!

# Feat: execute solidity

Let us begin with the most basic example. It is fine if you do not understand everything right now, we will go into more detail later. Create a file named "SimpleStorage.sol" with content:

```
pragma solidity ^0.4.24;

contract SimpleStorage {
    uint storedData;

    function set(uint x) public {
        storedData = x;
    }

    function get() view public returns (uint) {
        return storedData;
    }
}
```

Compile it:

```sh
$ solc --bin SimpleStorage.sol -o SimpleStorage
$ evm exec -code $(cat SimpleStorage/SimpleStorage.bin)
```

0x6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a7230582099c66a25d59f0aa78f7ebc40748fa1d1fbc335d8d780f284841b30e0365acd960029
