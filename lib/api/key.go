package api

import (
	"crypto/sha512"
	"errors"
	"net"
	"strings"
)

var encoding = []byte{
	99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 99, 99, 99, 99, 99, 99,
	99, 99, 10, 11, 12, 99, 13, 14, 15, 99, 16, 17, 18, 19, 20, 99,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 99, 99, 99, 99, 99,
	99, 99, 10, 11, 12, 99, 13, 14, 15, 99, 16, 17, 18, 19, 20, 99,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 99, 99, 99, 99, 99,
}

var ErrBadEncoding = errors.New("Bad encoding")

func KeyToAddr(key string) net.IP {
	if strings.HasSuffix(key, ".k") {
		key = strings.TrimSuffix(key, ".k")
		key = strings.ToLower(key)
		l := len(key)
		var d [32]byte
		idx := 0
		o_idx := 0
		var nextbyte, bits uint32
		for idx < l {
			c := key[idx]
			idx++
			if c&0x80 > 0 {
				return nil
			}
			b := encoding[c]
			if b > 31 {
				return nil
			}
			nextbyte |= uint32(b) << bits
			bits += 5
			if bits >= 8 {
				d[o_idx] = byte(nextbyte)
				o_idx++
				bits -= 8
				nextbyte >>= 8
			}
		}
		d1 := sha512.Sum512(d[:])
		d1 = sha512.Sum512(d1[:])
		return d1[:16]
	}
	return nil
}
