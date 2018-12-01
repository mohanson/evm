package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

const help = `usage: evm <command> [<args>]

The most commonly used daze commands are:
  disasm     Disassemble bytecode
  exec       Execute bytecode

Run 'evm <command> -h' for more information on a command.`

func printHelpAndExit() {
	fmt.Println(help)
	os.Exit(0)
}

func exDisasm() {
	var (
		flCode = flag.String("code", "0x603760005360005160005560016000f3", "bytecode")
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
}

func exExec() {
	var (
		flBlockNumber = flag.Int("number", 0, "block number")
		flCode        = flag.String("code", "0x603760005360005160005560016000f3", "bytecode")
		flCoinbase    = flag.String("coinbase", common.Address{}.String(), "coinbase")
		flData        = flag.String("data", "0x", "data")
		flDifficulty  = flag.Int("difficulty", 0, "difficulty")
		flGasLimit    = flag.Int("gaslimit", 100000, "gas limit")
		flGasPrice    = flag.Int("gasprice", 1, "gas price")
		flOrigin      = flag.String("origin", common.Address{}.String(), "sender")
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

	ret, sdb, err := runtime.Execute(common.FromHex(*flCode), common.FromHex(*flData), &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if sdb.Error() != nil {
		log.Fatalln(sdb.Error())
	}
	vm.WriteTrace(os.Stdout, slg.StructLogs())
	fmt.Println()
	fmt.Printf("Return = %#x\n", ret)
}

func main() {
	if len(os.Args) <= 1 {
		printHelpAndExit()
	}
	subCommand := os.Args[1]
	os.Args = os.Args[1:len(os.Args)]
	switch subCommand {
	case "disasm":
		exDisasm()
	case "exec":
		exExec()
	default:
		printHelpAndExit()
	}
}
