package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

var flVerbose = flag.Bool("v", false, "Be verbose")

func main() {
	flag.Parse()

	cfg := runtime.Config{}
	cfg.EVMConfig.Debug = true
	logger := vm.NewStructLogger(nil)
	cfg.EVMConfig.Tracer = logger

	code := common.FromHex("6060604052600a8060106000396000f360606040526008565b00")
	ret, _, err := runtime.Execute(code, nil, &cfg)
	if err != nil {
		fmt.Println(err)
	}
	gassum := 0
	for _, e := range logger.StructLogs() {
		gassum += int(e.GasCost)
		seps := []string{}
		seps = append(seps, fmt.Sprintf("[%04d]", e.Pc))
		if e.Op >= 0x60 && e.Op <= 0x7F {
			p1 := fmt.Sprintf("%s %#x", e.OpName(), code[e.Pc+1:e.Pc+uint64(e.Op)-0x5E])
			p2 := fmt.Sprintf("%-40s", p1)
			seps = append(seps, p2)
		} else {
			seps = append(seps, fmt.Sprintf("%-40s", e.OpName()))
		}
		seps = append(seps, "||")
		seps = append(seps, fmt.Sprintf("GasCost=%-5d Gas=%d", e.GasCost, e.Gas))
		fmt.Println(strings.Join(seps, " "))
		if *flVerbose {
			l := len(e.Stack)
			fmt.Printf("Stk.Len = %d\n", l)
			for i := 0; i < l; i++ {
				fmt.Printf("%04d: %#x\n", i, e.Stack[l-i-1])
			}
			fmt.Printf("Mem.Len = %d\n", e.MemorySize)

			for i := 0; i < e.MemorySize; i += 16 {
				src := e.Memory[i : i+16]
				dst := make([]string, 16)
				for i, e := range src {
					dst[i] = hex.EncodeToString([]byte{e})
				}
				fmt.Println(strings.Join(dst, " "))
			}
			fmt.Println()
		}
	}
	fmt.Println()
	fmt.Printf("Return: %#x\n", ret)
	fmt.Printf("Cost:   %d\n", gassum)
}
