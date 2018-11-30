package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

func test() {
	ret, _, err := runtime.Execute(common.FromHex("6060604052600a8060106000396000f360606040526008565b00"), nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
	// Output:
	// [96 96 96 64 82 96 8 86 91 0]
}

func main() {
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
	}
	fmt.Println()
	fmt.Printf("Return: %#x\n", ret)
	fmt.Printf("Cost: %d\n", gassum)
}
