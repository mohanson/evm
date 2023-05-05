package evm

import (
	"encoding/hex"
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
	dump := db.RawDump()
	for add, acc := range dump.Accounts {
		add := common.HexToAddress(add)
		for k := range acc.Storage {
			v := db.GetState(add, common.HexToHash(k))
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
		add := common.HexToAddress(add)
		db.CreateAccount(add)
		balance, _ := new(big.Int).SetString(acc.Balance, 16)
		db.SetBalance(add, balance)
		db.SetNonce(add, acc.Nonce)
		code, err := hex.DecodeString(acc.Code)
		if err != nil {
			return err
		}
		db.SetCode(add, code)
		for k, v := range acc.Storage {
			db.SetState(add, common.HexToHash(k), common.HexToHash(v))
		}
	}
	return nil
}
