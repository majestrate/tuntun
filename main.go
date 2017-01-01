package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/majestrate/tuntun/lib/api"
	"github.com/majestrate/tuntun/lib/api/admin"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func findFC00Addr() (laddr net.Addr, err error) {
	var ifs []net.Interface
	ifs, err = net.Interfaces()
	if err == nil {
		for _, netif := range ifs {
			var addrs []net.Addr
			addrs, err = netif.Addrs()
			if err == nil {
				for _, addr := range addrs {
					if strings.HasPrefix(addr.String(), "fc") {
						laddr = addr
						return
					}
				}
			}
		}
	}
	if err == nil {
		err = errors.New("cannot find fc00 address")
	}
	return
}

func main() {
	port := 1880
	adminfile := admin.DefaultAdminFile()
	laddr, e := findFC00Addr()
	if e == nil {
		addr := strings.Split(laddr.String(), "/")[0]
		addr = fmt.Sprintf("[%s]:%d", addr, port)
		log.Printf("serving on http://%s/", addr)

		handleNewRequest := func(w http.ResponseWriter, r *http.Request) {
			pubkey := r.URL.Query().Get("pubkey")
			pubaddr := api.KeyToAddr(pubkey)
			if pubaddr == nil {
				w.WriteHeader(400)
				io.WriteString(w, "bad pubkey")
				return
			}
			addr, _, _ := net.SplitHostPort(r.RemoteAddr)
			naddr := net.ParseIP(addr)
			if naddr.Equal(pubaddr) {
				a, err := admin.GetAdminFromFile(adminfile)
				if err == nil {
					s, err := a.Session()
					if err == nil {
						defer s.Close()
						info, err := s.AddTunnelIfNotThere(pubkey)
						if err == nil {
							json.NewEncoder(w).Encode(info)
							return
						}
					}
				}
				if err != nil {
					w.WriteHeader(500)
					io.WriteString(w, err.Error())
				}
			} else {
				w.WriteHeader(403)
			}
		}

		e = http.ListenAndServe(addr, http.HandlerFunc(handleNewRequest))
	}
	if e != nil {
		log.Fatal(e)
	}
}
