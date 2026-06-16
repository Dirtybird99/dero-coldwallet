# DERO ColdWallet

Generate a DERO **cold wallet** completely offline: an address, its 25-word
recovery seed, and a valid account-registration transaction. Designed to run on
an air-gapped machine so your seed is never created on a networked device.

Hardened community fork of [`8lecramm/dero-tools`](https://github.com/8lecramm/dero-tools)
(originally by Slixe). **DERO Stargate (DERO-HE) only.**

> Free, non-commercial, Research Use only. See `LICENSE` and `NOTICE.md`.

## Why this exists

DERO uses an account-based model: an address must be **registered on-chain
before it can receive funds**, and registration requires a small proof-of-work.
This tool computes that registration **offline**, so the air-gapped machine never
needs a network connection. You carry only the (non-secret) registration data
out to broadcast later.

## Security model

- **Your 25-word seed is the only thing that controls your funds.** The address
  and registration data are public.
- Every wallet is **self-verified**: the seed is re-derived back to the address
  before anything is printed, and the program aborts on any mismatch.
- The registration TX is validated with derohe's own `IsRegistrationValid()`.
- **No secret material is written to disk.** The seed lives only in memory unless
  you pass `--paper` (which writes it in plaintext). The program does create and then
  delete a throwaway, effectively-empty temporary data directory on startup (required
  by derohe's initialization); it is removed on normal exit, and left behind only
  if the process is terminated before it finishes (Ctrl+C / SIGINT, SIGTERM, or
  SIGKILL all skip cleanup) — in which case the leftover directory still holds no
  secret material.
- Runs in forced **offline mode**; it makes no network connections.

## Build

Requires the Go version recorded in `go.mod` (currently `go1.26.0`).

```
go build -o coldwallet ./cmd/coldwallet
```

For reproducible, distributable binaries, use `./build.sh` (see "Verify" below).

## Usage

```
# generate a mainnet wallet (English seed)
./coldwallet --language=0

# list seed languages
./coldwallet --list-languages

# skip the registration proof-of-work (faster; register before funding)
./coldwallet --language=0 --no-register

# testnet
./coldwallet --language=0 --testnet

# also write a printable, fully-offline HTML paper wallet
./coldwallet --language=0 --paper=wallet.html
```

The `--paper` file contains the **seed in plaintext**. Print it, then securely
delete the file. It is self-contained (no network requests) and prints on
A4/Letter, including QR codes for the address, seed, and registration.

## Recommended air-gapped procedure

1. Boot a verified live OS (for example Tails) on a machine with networking
   physically disabled. Use wired peripherals.
2. Build or copy a verified `coldwallet` binary (see "Verify").
3. Run it. Record the **seed onto metal**, stored in multiple secure locations.
   Never photograph it or store it on any online device.
4. Carry only the **registration hex** (not secret) to an online machine and
   broadcast it with the `curl` command the tool prints (replace the node
   address). Confirm the address is registered on a DERO explorer before funding.
5. Recover any time with the official DERO wallet using the 25 words.

## Verify before you run

This program is built reproducibly. You do not have to trust a binary we
produced; you can confirm it byte-for-byte from source.

**Rebuild and compare:**

```
git checkout <release-tag>
./build.sh
# compare your dist/SHA256SUMS against the published SHA256SUMS
```

`build.sh` pins every reproducibility input — `CGO_ENABLED=0`, `-trimpath`,
`-ldflags=all=-buildid= -s -w`, `-buildvcs=false`, and the microarch baseline
(`GOAMD64=v1` / `GOARM64=v8.0`) — so the same source produces byte-identical
binaries **on the pinned Go toolchain** (`toolchain go1.26.0` in `go.mod`; with
the default `GOTOOLCHAIN=auto`, `go` fetches it automatically), regardless of
whether you obtained the source via `git clone` or a release tarball.
(`-buildvcs=false` matters: without it the toolchain stamps the `.git` revision
and dirty-state into the binary, so clone-based and tarball-based rebuilders
would otherwise never converge.) **Run `./build.sh` as-is** rather than
reproducing these flags by hand — a partial flag set will not match the published
hashes. A different Go version may also produce a different binary, so match the
toolchain when comparing. Dependencies are pinned by `go.sum`.

**Check a downloaded release:**

```
sha256sum -c SHA256SUMS
minisign -Vm SHA256SUMS -P <published-public-key>
```

Independent rebuilders are encouraged to add their matching hashes to `SIGS.md`.

## Independent verification helper

`cmd/verify` re-derives an address from a seed and validates a registration TX
offline, using the official derohe code paths:

```
go run ./cmd/verify --seed="..." --address=dero1... --tx=<hex>
```
