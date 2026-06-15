module github.com/8lecramm/dero-tools

go 1.25.0

// Pinned for reproducible builds: rebuilders get byte-identical binaries only
// on this exact toolchain. With GOTOOLCHAIN=auto (the default) `go` fetches it.
toolchain go1.26.0

require (
	github.com/deroproject/derohe v0.0.0-20250813215012-9b6a8b82c839
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/makiuchi-d/gozxing v0.1.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
)

require (
	github.com/VictoriaMetrics/metrics v1.43.2 // indirect
	github.com/beevik/ntp v1.5.0 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/caarlos0/env/v6 v6.10.1 // indirect
	github.com/cenkalti/rpc2 v1.0.5 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/coder/websocket v1.8.15 // indirect
	github.com/creachadair/jrpc2 v1.3.5 // indirect
	github.com/creachadair/mds v0.26.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/deroproject/graviton v0.0.0-20220130070622-2c248a53b2e1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fxamacker/cbor/v2 v2.9.2 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/klauspost/reedsolomon v1.14.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lesismal/llib v1.2.2 // indirect
	github.com/lesismal/nbio v1.6.9 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xtaci/kcp-go/v5 v5.6.72 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
