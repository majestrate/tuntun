package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func RandStr(l int) string {
	b := make([]byte, l)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)[:l]
}
