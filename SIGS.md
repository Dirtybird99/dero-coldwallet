# Reproducible-build attestations

This file collects independent confirmations that the published release binaries
match what the source produces. The more independent rebuilders agree, the less
anyone has to trust a single maintainer or machine (the model used by Bitcoin
Core's `guix.sigs` and Monero).

## How to attest

1. `git checkout <release-tag>`
2. `./build.sh`
3. Confirm your `dist/SHA256SUMS` matches the released `SHA256SUMS`.
4. Add a row below with your name/handle, the release tag, your Go version, and
   "match", then open a pull request.

## Attestations

| Release | Rebuilder | Go version | Result |
|---------|-----------|------------|--------|
| (none yet) | | | |
