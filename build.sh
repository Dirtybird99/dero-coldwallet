#!/usr/bin/env bash
# Reproducible build for DERO ColdWallet.
#
# Pure Go, no cgo, so byte-for-byte reproducibility comes from these inputs:
#   CGO_ENABLED=0          remove the host C toolchain as an input
#   -trimpath              remove local filesystem paths from the binary
#   -ldflags=all=-buildid= strip the non-deterministic build id
#   -buildvcs=false        do NOT stamp .git revision/time/dirty state into the
#                          binary; otherwise clone vs. tarball rebuilders diverge
#   GOAMD64=v1 / GOARM64=v8.0  pin the microarch baseline so a rebuilder whose
#                          environment overrides it still converges
# Dependencies are pinned by go.sum; the Go toolchain is pinned by go.mod
# (asserted below against the actual toolchain in use).
#
# Anyone can run this on the same Go version and get identical binaries,
# then compare against the published SHA256SUMS. No need to trust our binary.
set -euo pipefail

export CGO_ENABLED=0
export GOFLAGS=-mod=readonly
export GOAMD64=v1
export GOARM64=v8.0
LDFLAGS="all=-buildid= -s -w"
OUT="dist"

GO_VERSION="$(go env GOVERSION)"
EXPECT="$(awk '/^toolchain /{print $2}' go.mod)"
if [ -n "$EXPECT" ] && [ "$GO_VERSION" != "$EXPECT" ]; then
  echo "Toolchain mismatch: go.mod pins $EXPECT but building with $GO_VERSION" >&2
  echo "With the default GOTOOLCHAIN=auto, 'go' fetches the pinned toolchain;" >&2
  echo "if you forced GOTOOLCHAIN=local, install/select $EXPECT and retry." >&2
  exit 1
fi
echo "Building with $GO_VERSION (reproducible flags)"

rm -rf "$OUT"
mkdir -p "$OUT"

targets=(
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
  "darwin/amd64"
  "darwin/arm64"
)

for t in "${targets[@]}"; do
  os="${t%/*}"; arch="${t#*/}"
  ext=""; [ "$os" = "windows" ] && ext=".exe"
  out="$OUT/coldwallet-${os}-${arch}${ext}"
  echo "  $t -> $out"
  GOOS="$os" GOARCH="$arch" go build -trimpath -buildvcs=false -ldflags="$LDFLAGS" -o "$out" ./cmd/coldwallet
done

( cd "$OUT" && sha256sum coldwallet-* > SHA256SUMS )
echo
echo "Artifacts and checksums written to $OUT/"
cat "$OUT/SHA256SUMS"
echo
echo "To sign (maintainer): minisign -Sm $OUT/SHA256SUMS"
echo "To verify (user):     minisign -Vm $OUT/SHA256SUMS -P <public-key>  &&  sha256sum -c $OUT/SHA256SUMS"
