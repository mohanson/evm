# EVM

Inspired by <https://github.com/ethereum/go-ethereum/tree/master/cmd/evm>, though `evm` has been modified to support persistence.

# Installation

```sh
$ go get github.com/mohanson/evm/cmd/evm
```

# Feat: disassemble bytecode

Here is some very simple bytecode I wrote:

```text
0x6005600401
```

To disassemble, run `evm disasm -code 0x6005600401`, which produces:

```text
[0000] PUSH1 0x05
[0002] PUSH1 0x04
[0004] ADD
```

# Feat: execute bytecode

Let's use a test case to demonstrate: [return0.json](https://github.com/ethereum/tests/blob/develop/VMTests/vmSystemOperations/return0.json). To execute the test, run

```sh
$ evm exec -code 0x603760005360005160005560016000f3
```

```text
[Many Outputs]
...

Return  = 0x37
```

# Feat: execute solidity

Let us begin with the most basic example. It is fine if you do not understand everything right now, we will go into more detail later. Create a file named "SimpleStorage.sol" with content:

```text
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

Use `solc` to compile it to bytecode:

```sh
$ solc --bin SimpleStorage.sol

# The output bytecode. I will use the $BYTECODE to instead of the original text.
608060405234801561001057600080fd5b5060df8061001f6000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a7230582099c66a25d59f0aa78f7ebc40748fa1d1fbc335d8d780f284841b30e0365acd960029
```

```sh
# create it at /tmp/db.json.
$ evm create -db /tmp/db.json
             -data $BYTECODE

Cost    = 55307
Address = 0xBd770416a3345F91E4B34576cb804a576fa48EB1

# let's call SimpleStorage.set(42)
$ evm call -db /tmp/db.json
           -address 0xBd770416a3345F91E4B34576cb804a576fa48EB1
           -data 0x60fe47b1000000000000000000000000000000000000000000000000000000000000002a

# let's call SimpleStorage.get()
$ evm call -db /tmp/db.json
           -address 0xBd770416a3345F91E4B34576cb804a576fa48EB1
           -data 0x6d4ce63c

Cost    = 424
Return  = 000000000000000000000000000000000000000000000000000000000000002a
```

Perfect!

---

Note: data `0x60fe47b1000000000000000000000000000000000000000000000000000000000000002a` = `0x60fe47b1` + `0x2a`, where `60fe47b1` provided by `python -c "import eth_utils; print(eth_utils.keccak(b'set(uint256)')[0:4].hex())"`.
