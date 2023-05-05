package main

import (
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/mohanson/evm"
)

const help = `usage: evm <command> [<args>]

The most commonly used daze commands are:
  disasm     Disassemble bytecode
  exec       Execute bytecode
  create     Create a contract
  call       Call contract

Run 'evm <command> -h' for more information on a command.`

func printHelpAndExit() {
	fmt.Println(help)
	os.Exit(0)
}

func exDisasm() error {
	var (
		flCode = flag.String("code", "", "bytecode")
	)
	flag.Parse()
	code := common.FromHex(*flCode)
	for pc := 0; pc < len(code); pc++ {
		op := vm.OpCode(code[pc])
		fmt.Printf("[%04d] %v", pc, op)
		e := int(op)
		if e >= 0x60 && e <= 0x7F {
			l := e - int(vm.PUSH1) + 1
			off := pc + 1
			end := func() int {
				if len(code) < off+l {
					return len(code)
				}
				return off + l
			}()
			data := make([]byte, l)
			copy(data, code[off:end])
			fmt.Printf(" %#x", data)
			pc += l
		}
		fmt.Println()
	}
	return nil
}

func exMacall(subcmd string) error {
	var (
		flAddress     = flag.String("address", "", "address")
		flBlockNumber = flag.Int("number", 0, "block number")
		flCode        = flag.String("code", "", "bytecode")
		flCoinbase    = flag.String("coinbase", "", "coinbase")
		flData        = flag.String("data", "", "data")
		flDB          = flag.String("db", "", "database")
		flDifficulty  = flag.Int("difficulty", 0, "difficulty")
		flGasLimit    = flag.Int("gaslimit", 100000, "gas limit")
		flGasPrice    = flag.Int("gasprice", 1, "gas price")
		flOrigin      = flag.String("origin", "", "sender")
		flValue       = flag.Int64("value", 0, "value")
	)
	flag.Parse()
	cfg := runtime.Config{}
	cfg.BlockNumber = big.NewInt(int64(*flBlockNumber))
	cfg.Coinbase = common.HexToAddress(*flCoinbase)
	cfg.Difficulty = big.NewInt(int64(*flDifficulty))
	cfg.GasLimit = uint64(*flGasLimit)
	cfg.GasPrice = big.NewInt(int64(*flGasPrice))
	cfg.Origin = common.HexToAddress(*flOrigin)
	cfg.Value = big.NewInt(*flValue)
	cfg.EVMConfig.Debug = true
	slg := vm.NewStructLogger(nil)
	cfg.EVMConfig.Tracer = slg
	sdb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	if err != nil {
		return err
	}
	if subcmd == "create" || subcmd == "call" {
		if *flDB == "" {
			return errors.New("evm: missing -db operand")
		}
		if err := evm.LoadStateDB(sdb, *flDB); err != nil {
			if os.IsExist(err) {
				return err
			}
		}
	}
	cfg.State = sdb
	switch subcmd {
	case "exec":
		ret, _, err := runtime.Execute(common.FromHex(*flCode), common.FromHex(*flData), &cfg)
		if err != nil {
			return err
		}
		vm.WriteTrace(os.Stdout, slg.StructLogs())
		fmt.Println()
		fmt.Println("Return  =", common.Bytes2Hex(ret))
		return nil
	case "create":
		_, add, gas, err := runtime.Create(common.FromHex(*flData), &cfg)
		if err != nil {
			return err
		}
		vm.WriteTrace(os.Stdout, slg.StructLogs())
		fmt.Println()
		fmt.Println("Cost    =", *flGasLimit-int(gas))
		fmt.Println("Address =", add.String())
		return evm.SaveStateDB(sdb, *flDB)
	case "call":
		ret, gas, err := runtime.Call(common.HexToAddress(*flAddress), common.FromHex(*flData), &cfg)
		if err != nil {
			return err
		}
		vm.WriteTrace(os.Stdout, slg.StructLogs())
		fmt.Println()
		fmt.Println("Cost    =", *flGasLimit-int(gas))
		fmt.Println("Return  =", common.Bytes2Hex(ret))
		return evm.SaveStateDB(sdb, *flDB)
	}
	return nil
}

func main() {
	if len(os.Args) <= 1 {
		printHelpAndExit()
	}
	subCommand := os.Args[1]
	os.Args = os.Args[1:len(os.Args)]
	var err error
	switch subCommand {
	case "disasm":
		err = exDisasm()
	case "exec", "create", "call":
		err = exMacall(subCommand)
	default:
		printHelpAndExit()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
