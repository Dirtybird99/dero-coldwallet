// Command verify independently checks a generated DERO cold wallet:
//  1. seed -> address round-trip (does the seed actually control the address?)
//  2. the registration TX is cryptographically valid (IsRegistrationValid) and
//     binds to the same public key -- proven offline, no broadcast required.
//
// This uses the official derohe walletapi/transaction code paths -- the same
// ones dero-wallet-cli uses on recovery.
package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/deroproject/derohe/transaction"
	"github.com/deroproject/derohe/walletapi"
)

func main() {
	seed := flag.String("seed", "", "25-word recovery seed")
	expected := flag.String("address", "", "expected DERO mainnet address (dero1...)")
	txhex := flag.String("tx", "", "registration TX hex")
	flag.Parse()

	if *seed == "" || *expected == "" {
		fmt.Println("usage: verify --seed=\"...\" --address=dero1... [--tx=<hex>]")
		os.Exit(2)
	}

	account, err := walletapi.Generate_Account_From_Recovery_Words(*seed)
	if err != nil {
		fmt.Println("FAIL: seed did not parse:", err)
		os.Exit(1)
	}
	addr := account.GetAddress()
	addr.Mainnet = true
	got := addr.String()

	fmt.Println("Regenerated address:", got)
	fmt.Println("Expected   address:", *expected)
	if got != *expected {
		fmt.Println("RESULT: ❌ MISMATCH — the seed does NOT control this address. DO NOT FUND.")
		os.Exit(1)
	}
	fmt.Println("RESULT: ✅ seed ↔ address round-trip OK (the seed controls the address).")

	if *txhex != "" {
		raw, err := hex.DecodeString(*txhex)
		if err != nil {
			fmt.Println("FAIL: tx hex did not decode:", err)
			os.Exit(1)
		}
		var tx transaction.Transaction
		if err := tx.Deserialize(raw); err != nil {
			fmt.Println("FAIL: tx did not deserialize:", err)
			os.Exit(1)
		}
		fmt.Println("\nTX type:", tx.TransactionType.String())
		if tx.TransactionType != transaction.REGISTRATION {
			fmt.Println("RESULT: ❌ not a REGISTRATION transaction.")
			os.Exit(1)
		}
		// the registration must bind to THIS wallet's public key
		pub := account.Keys.Public.EncodeCompressed()
		if !bytes.Equal(tx.MinerAddress[:], pub[:]) {
			fmt.Println("RESULT: ❌ registration public key does not match the wallet.")
			os.Exit(1)
		}
		if !tx.IsRegistrationValid() {
			fmt.Println("RESULT: ❌ registration signature INVALID.")
			os.Exit(1)
		}
		fmt.Println("RESULT: ✅ registration TX is valid and binds to this wallet's key (verified offline).")
	}
}
