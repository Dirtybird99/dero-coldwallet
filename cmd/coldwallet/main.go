// DERO ColdWallet — offline wallet + registration generator for DERO Stargate (DERO-HE).
//
// Originally by Slixe, modified by mmarcel; hardened community fork.
//
// It generates, fully offline, a DERO address, its 25-word recovery seed, and a
// valid account-registration transaction. Before printing anything it re-derives
// the address from the seed and aborts on any mismatch, and it validates the
// registration transaction with derohe's own IsRegistrationValid(). Output goes
// to stdout only — nothing is written to disk.
package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
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

var cmd = `DERO ColdWallet — offline wallet + registration generator (Stargate/DERO-HE)
Originally by Slixe, modified by mmarcel; hardened community fork.

Generates a DERO address, its 25-word seed and a valid registration TX fully
offline. The seed is re-derived back to the address and the registration TX is
validated before printing. Output is printed to stdout only — nothing is written
to disk. Run this on an air-gapped machine.

Usage:
  coldwallet --language=<id> [--testnet] [--no-register] [--paper=<path>]
  coldwallet --list-languages
  coldwallet -h | --help

Options:
  --language=<id>     Seed language id (run --list-languages to see the choices)
  --list-languages    List the available seed languages and exit
  --testnet           Generate a testnet address (default: mainnet)
  --no-register       Skip the registration proof-of-work (faster; the address
                      cannot receive funds until it is registered later)
  --paper=<path>      Also write a self-contained, offline HTML paper wallet to
                      <path>. WARNING: this file contains the SEED in plaintext.
                      Print it, then securely delete the file.
  -h --help           Show this help
`

var language mnemonics.Language

func main() {
	os.Exit(run())
}

func run() int {
	arguments, err := docopt.Parse(cmd, nil, true, "DERO ColdWallet", false)
	if err != nil {
		fmt.Println("Error while parsing options:", err)
		return 2
	}

	testnet := arguments["--testnet"] != nil && arguments["--testnet"].(bool)
	mainnet := !testnet

	// Keep zero filesystem footprint: derohe's globals.Initialize() does a
	// MkdirAll on the data directory (defaults to the CWD). Point it at a
	// throwaway temp dir and remove it on exit so nothing is left behind.
	tmpDir, err := os.MkdirTemp("", "coldwallet-")
	if err != nil {
		fmt.Println("Error creating temp dir:", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	globals.Arguments = arguments
	globals.Arguments["--debug"] = false
	globals.Arguments["--testnet"] = testnet
	globals.Arguments["--data-dir"] = tmpDir
	globals.Initialize()

	if v := arguments["--list-languages"]; v != nil && v.(bool) {
		for i, l := range mnemonics.Language_List() {
			fmt.Printf("(%d): %s\n", i, l)
		}
		return 0
	}

	if value := arguments["--language"]; value != nil {
		if s, err := strconv.Atoi(value.(string)); err == nil && s >= 0 && s < len(mnemonics.Languages) {
			language = mnemonics.Languages[s]
		} else {
			fmt.Println("Invalid language id. Run 'coldwallet --list-languages' to see the choices.")
			return 2
		}
	} else {
		fmt.Println("No seed language selected. Run 'coldwallet --list-languages' to see the choices.")
		return 2
	}

	net := "MAINNET"
	if testnet {
		net = "TESTNET"
	}
	fmt.Println("DERO ColdWallet —", config.Version, "—", net)
	fmt.Println("Seed language:", language.Name)
	fmt.Println("Generating new random wallet (offline)...")

	w, err := walletapi.Create_Encrypted_Wallet_Random_Memory("")
	if err != nil {
		fmt.Println("Error while generating random wallet:", err)
		return 1
	}
	w.SetOfflineMode()
	w.SetNetwork(mainnet)

	address := w.GetAddress().String()
	seed := w.GetSeedinLanguage(language.Name)

	// --- self-verification: the seed MUST re-derive the displayed address ---
	if err := verifySeedMatchesAddress(seed, address, mainnet); err != nil {
		fmt.Println("\n❌ FATAL: seed/address self-check FAILED —", err)
		fmt.Println("Do NOT use this wallet. This should never happen; please report it.")
		return 1
	}

	fmt.Println()
	fmt.Println("Address:", address)
	fmt.Println("Seed   :", seed)
	fmt.Println("(seed ↔ address self-check: OK)")

	var txData string
	register := !(arguments["--no-register"] != nil && arguments["--no-register"].(bool))
	if !register {
		fmt.Println("\nRegistration skipped (--no-register). Register later before funding.")
	} else {
		fmt.Println("\nGenerating valid registration TX (proof-of-work, may take a minute)...")
		regTx := mineRegistration(w)

		// validate the registration with derohe's own check before trusting it
		if regTx.TransactionType != transaction.REGISTRATION || !regTx.IsRegistrationValid() {
			fmt.Println("❌ FATAL: produced registration TX failed IsRegistrationValid(). Aborting.")
			return 1
		}

		txData = hex.EncodeToString(regTx.Serialize())
		port := config.Mainnet.RPC_Default_Port
		if testnet {
			port = config.Testnet.RPC_Default_Port
		}
		fmt.Println("TX data:", txData)
		fmt.Println("(registration validity check: OK)")
		fmt.Printf("\nBroadcast it from an online machine (replace <node_address>):\n")
		fmt.Printf("curl -X POST <node_address>:%d/json_rpc -H \"Content-Type: application/json\" -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"sendrawtransaction\",\"params\":{\"tx_as_hex\":\"%s\"}}'\n", port, txData)
	}

	if value := arguments["--paper"]; value != nil {
		path := value.(string)
		if err := WritePaperWallet(path, net, config.Version.String(), address, seed, txData); err != nil {
			fmt.Println("\nError writing paper wallet:", err)
			return 1
		}
		fmt.Println("\n📄 Paper wallet written to:", path)
		fmt.Println("   WARNING: that file contains the SEED in plaintext. Print it, then securely delete it.")
	}

	printSecurityNotice(seed)
	return 0
}

// verifySeedMatchesAddress re-derives the address from the seed using the
// official recovery path (the same one dero-wallet-cli uses) and compares it.
func verifySeedMatchesAddress(seed, expected string, mainnet bool) error {
	account, err := walletapi.Generate_Account_From_Recovery_Words(seed)
	if err != nil {
		return fmt.Errorf("seed did not parse: %w", err)
	}
	addr := account.GetAddress()
	addr.Mainnet = mainnet
	if got := addr.String(); got != expected {
		return fmt.Errorf("re-derived %s != %s", got, expected)
	}
	return nil
}

func printSecurityNotice(seed string) {
	fmt.Println("\n--- SECURITY ---")
	fmt.Println("• The SEED above is the ONLY thing that controls your funds. Back it up")
	fmt.Println("  on paper or metal, in multiple safe locations. Never store it digitally.")
	fmt.Println("• The TX data and address are NOT secret. Only the seed is.")
	fmt.Println("• Generate on an air-gapped machine. Only carry the TX data or address out.")
	_ = seed
}

// mineRegistration runs the registration proof-of-work across all cores and
// returns the first valid transaction found.
func mineRegistration(w *walletapi.Wallet_Memory) *transaction.Transaction {
	txChan := make(chan *transaction.Transaction, 1)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go GetRegistrationTX(w, txChan)
	}
	regTx := <-txChan
	return regTx
}

// faster version: incrementing the secret and point instead of generating random
// secrets and points until we find a valid one. Safe because exactly one
// registration is ever published per wallet (no cross-message nonce reuse).
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

	select {
	case txChan <- &tx:
	default: // another goroutine already delivered a result
	}
}
