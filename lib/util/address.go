package util

import (
	"net"
)

type AddressGenerator interface {
	GetIPFor(pubkey string) (net.IP, error)
}
