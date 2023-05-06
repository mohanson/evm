package evm

import (
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

func SaveStateDB(db *state.StateDB, fn string) error {
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	db.Commit(false)
	db.IntermediateRoot(true)
	dump := db.RawDump(nil)
	for add, acc := range dump.Accounts {
		for k := range acc.Storage {
			v := db.GetState(add, k)
			acc.Storage[k] = v.String()
		}
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(dump)
}

func LoadStateDB(db *state.StateDB, fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	dump := state.Dump{}
	if err := json.NewDecoder(f).Decode(&dump); err != nil {
		return err
	}
	for add, acc := range dump.Accounts {
		db.CreateAccount(add)
		balance, _ := new(big.Int).SetString(acc.Balance, 16)
		db.SetBalance(add, balance)
		db.SetNonce(add, acc.Nonce)
		db.SetCode(add, acc.Code)
		for k, v := range acc.Storage {
			db.SetState(add, k, common.HexToHash(v))
		}
	}
	return nil
}
