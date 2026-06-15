package main

import "testing"

// Burned test vector: this wallet was generated during development and its
// registration was broadcast to mainnet, so the seed is intentionally public.
// Never use it for funds. It exists only to pin the seed -> address derivation.
const (
	testSeed = "anybody auburn down awning nightly eels icon ungainly jump upstairs foolish swung rudely kangaroo catch gossip upcoming kennel echo dwindling bifocals potato jewels jailed bifocals"
	testAddr = "dero1qyfez0fm768fmp9tele8crqnvq59jgmjcx07y85y9x8mv0a43fss2qgx3n4pl"
)

func TestVerifySeedMatchesAddress_RoundTrip(t *testing.T) {
	if err := verifySeedMatchesAddress(testSeed, testAddr, true); err != nil {
		t.Fatalf("expected seed to match address, got error: %v", err)
	}
}

func TestVerifySeedMatchesAddress_WrongAddress(t *testing.T) {
	wrong := "dero1qyfez0fm768fmp9tele8crqnvq59jgmjcx07y85y9x8mv0a43fss2qgx3n4XX"
	if err := verifySeedMatchesAddress(testSeed, wrong, true); err == nil {
		t.Fatal("expected mismatch error for a wrong address, got nil")
	}
}

func TestVerifySeedMatchesAddress_TestnetDiffersFromMainnet(t *testing.T) {
	// the same seed must NOT match the mainnet address when rendered as testnet
	if err := verifySeedMatchesAddress(testSeed, testAddr, false); err == nil {
		t.Fatal("expected testnet rendering to differ from mainnet address, got nil")
	}
}

func TestVerifySeedMatchesAddress_BadSeed(t *testing.T) {
	if err := verifySeedMatchesAddress("not a real seed", testAddr, true); err == nil {
		t.Fatal("expected parse error for an invalid seed, got nil")
	}
}
