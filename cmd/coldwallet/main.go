package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"runtime"
	"strconv"

	"github.com/deroproject/derohe/config"
	"github.com/deroproject/derohe/cryptography/bn256"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/transaction"
	"github.com/deroproject/derohe/walletapi"
	"github.com/deroproject/derohe/walletapi/mnemonics"
	"github.com/docopt/docopt-go"
)

var cmd = `DERO ColdWallet by Slixe, modified by mmarcel
DERO : A secure, private blockchain with smart-contracts

Usage:
  coldwallet [options] 
  coldwallet -h | --help

Options:
  --languages						Show languages available
  --seed=<seed>						Seed
  --language=<id>					Seed language
`
var language mnemonics.Language

func main() {
	arguments, err := docopt.Parse(cmd, nil, true, "DERO ColdWallet", false)
	if err != nil {
		fmt.Println("Error while parsing options err:", err)
	}

	globals.Arguments = arguments
	globals.Arguments["--debug"] = false
	globals.Arguments["--testnet"] = false
	globals.Initialize()

	if value := globals.Arguments["--languages"]; value != nil {
		for i, l := range mnemonics.Language_List() {
			fmt.Printf("(%d): %s\n", i, l)
		}
		return
	}

	if value := globals.Arguments["--language"]; value != nil {
		if s, err := strconv.Atoi(value.(string)); err == nil && s >= 0 && s < len(mnemonics.Languages) {
			language = mnemonics.Languages[s]
			fmt.Println("Language seed selected:", language.Name)
		} else {
			fmt.Println("Invalid choice, please select a valid language.")
			return
		}
	} else {
		fmt.Println("No language for seed selected!")
		fmt.Println("Available:")
		for i, l := range mnemonics.Language_List() {
			fmt.Printf("(%d): %s\n", i, l)
		}
		return
	}

	fmt.Println("Version:", config.Version)
	fmt.Println("Generating new random wallet...")
	w, err := walletapi.Create_Encrypted_Wallet_Random_Memory("")
	if err != nil {
		fmt.Println("Error while generating random wallet:", err)
		return
	}
	w.SetOfflineMode()
	w.SetNetwork(true)

	fmt.Println("Address:", w.GetAddress())
	fmt.Println("Seed   :", w.GetSeedinLanguage(language.Name))

	txChan := make(chan *transaction.Transaction, 1)

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func() {
			GetRegistrationTX(w, txChan)
		}()
	}
	regTx := <-txChan
	close(txChan)
	tx_data := hex.EncodeToString(regTx.Serialize())
	fmt.Println("TX data:", tx_data)
	fmt.Printf("\ncurl -X POST <node_address>:10102/json_rpc -H \"Content-Type: application/json\" -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"sendrawtransaction\",\"params\":{\"tx_as_hex\":\"%s\"}}'\n", tx_data)
}

// faster version: incrementing the secret and point instead of generating random secrets and points until we find a valid one
func GetRegistrationTX(w *walletapi.Wallet_Memory, txChan chan<- *transaction.Transaction) {
	var tx transaction.Transaction
	tx.Version = 1
	tx.TransactionType = transaction.REGISTRATION
	add := w.GetAddress().PublicKey.EncodeCompressed()
	copy(tx.MinerAddress[:], add[:])

	var tmppoint bn256.G1

	tmpsecret := crypto.RandomScalar()
	tmppoint.ScalarMult(crypto.G, tmpsecret)

	for {
		serialize := []byte(fmt.Sprintf("%s%s", w.Get_Keys().Public.G1().String(), tmppoint.String()))
		c := crypto.ReducedHash(serialize)
		s := new(big.Int).Mul(c, w.Get_Keys().Secret.BigInt())
		s = s.Mod(s, bn256.Order)
		s = s.Add(s, tmpsecret)
		s = s.Mod(s, bn256.Order)

		crypto.FillBytes(c, tx.C[:])
		crypto.FillBytes(s, tx.S[:])

		hash := tx.GetHash()
		if hash[0] == 0 && hash[1] == 0 && hash[2] == 0 {
			break
		}
		tmpsecret.Add(tmpsecret, big.NewInt(1))
		tmppoint.Add(&tmppoint, crypto.G)
	}

	txChan <- &tx
}
