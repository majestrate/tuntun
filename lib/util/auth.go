package util

import "net"

type ClientAuthPolicy interface {
	Allow(ip net.IP) bool
}

type NullAuthPolicy struct {
}

func (p *NullAuthPolicy) Allow(ip net.IP) bool {
	return true
}
