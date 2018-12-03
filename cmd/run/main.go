package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/mohanson/evm"
)

func main() {
	main_1()
}

func main_0() {
	account := common.HexToAddress("000000000000000000000000636f6e7472616374")
	k := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	v := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000080")

	sdb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	if err != nil {
		log.Fatalln(err)
	}
	sdb.SetState(account, k, v)
	sdb.Commit(false)
	log.Println(sdb.GetState(account, k).String())
}

func main_1() {
	account := common.HexToAddress("000000000000000000000000636f6e7472616374")
	k := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	// v := common.HexToHash("000000000000000000000000000000000000000000000000000000000000002a")

	sdb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	if err != nil {
		log.Fatalln(err)
	}
	if err := evm.LoadStateDB(sdb, "/tmp/db.json"); err != nil {
		log.Fatalln(err)
	}
	log.Println(sdb.GetState(account, k))
}

func pick() {
	_ = log.Fatal
	_ = fmt.Errorf
	_ = big.NewInt
	_ = evm.SaveStateDB
}
