package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const help = `usage: evmfun <command> [<args>]

The most commonly used daze commands are:
  disasm     Disassemble bytecode

Run 'evmfun <command> -h' for more information on a command.`

func printHelpAndExit() {
	fmt.Println(help)
	os.Exit(0)
}

func exDisasm() {
	codeRaw, _ := ioutil.ReadAll(os.Stdin)
	codeStr := string(codeRaw)
	codeStr = strings.TrimSpace(codeStr)
	code := common.FromHex(codeStr)
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

func main() {
	if len(os.Args) <= 1 {
		printHelpAndExit()
	}
	subCommand := os.Args[1]
	os.Args = os.Args[1:len(os.Args)]
	switch subCommand {
	case "disasm":
		exDisasm()
	default:
		printHelpAndExit()
	}
}
