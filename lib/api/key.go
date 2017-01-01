package api

import (
	"crypto/sha512"
	"encoding/base32"
	"net"
	"strings"
)

func KeyToAddr(key string) net.IP {
	if strings.HasSuffix(key, ".k") {
		l := len(key)
		d, err := base32.StdEncoding.DecodeString(key[:l-2])
		if err == nil {
			d1 := sha512.Sum512(d)
			d1 = sha512.Sum512(d1[:])
			return d1[:16]
		}
	}
	return nil
}
