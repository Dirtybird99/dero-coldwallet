package main

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/png"
	"strings"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// decodeQR turns a data: URI produced by qrDataURI back into its text, the way
// a phone camera or QR reader would. This guards the recovery path: a QR that
// renders but decodes to truncated/mangled data would be silent fund loss.
func decodeQR(t *testing.T, dataURI string) string {
	t.Helper()
	const prefix = "data:image/png;base64,"
	if !strings.HasPrefix(dataURI, prefix) {
		t.Fatalf("unexpected data URI prefix: %.30q", dataURI)
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(dataURI, prefix))
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("png decode: %v", err)
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		t.Fatalf("bitmap: %v", err)
	}
	res, err := qrcode.NewQRCodeReader().Decode(bmp, nil)
	if err != nil {
		t.Fatalf("qr decode: %v", err)
	}
	return res.GetText()
}

func TestQRRoundTrip(t *testing.T) {
	cases := map[string]string{
		"address": "dero1qyfez0fm768fmp9tele8crqnvq59jgmjcx07y85y9x8mv0a43fss2qgx3n4pl",
		"seed":    testSeed,
		"tx":      "0100000113913d3bf68e9d84abcff27c0c136028592372c19fe21e84298fb63fb58a61050111dfa2948c69d1ffed09a9600e23a8e39b84d9392c6d0e01844606d60cb4f01f101962776c0942de7c47efcf3cdf87abafdb274b7034f88f73838a10b723ee97",
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			uri, err := qrDataURI(content, 320)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			got := decodeQR(t, string(uri))
			if got != content {
				t.Fatalf("QR round-trip mismatch:\n got: %q\nwant: %q", got, content)
			}
		})
	}
}
