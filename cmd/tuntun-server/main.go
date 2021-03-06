package main

import (
	"encoding/json"
	"fmt"
	"github.com/majestrate/tuntun/lib/api"
	"github.com/majestrate/tuntun/lib/api/admin"
	"github.com/majestrate/tuntun/lib/config"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	fname := "tuntun.ini"
	if len(os.Args) > 1 {
		fname = os.Args[1]
	}
	conf, err := config.Load(fname)

	if err != nil {
		log.Fatal(err)
	}

	adminfile := admin.DefaultAdminFile()
	a, e := admin.GetAdminFromFile(adminfile)
	if e != nil {
		log.Fatal(e)
	}
	s, e := a.Session()
	if e != nil {
		log.Fatal(e)
	}
	pk, e := s.GetOurPubkey()
	ip := api.KeyToAddr(pk)
	s.Close()
	gen := conf.Exit.AddressGenerator()
	auth := conf.Auth.AuthPolicy()
	if e == nil {
		addr := fmt.Sprintf("[%s]:%d", ip.String(), conf.Port)
		log.Printf("serving on http://%s/", addr)

		handleNewRequest := func(w http.ResponseWriter, r *http.Request) {
			pubkey := r.URL.Query().Get("pubkey")
			pubaddr := api.KeyToAddr(pubkey)
			if pubaddr == nil {
				w.WriteHeader(400)
				io.WriteString(w, "bad pubkey, "+pubkey)
				return
			}
			addr, _, _ := net.SplitHostPort(r.RemoteAddr)
			naddr := net.ParseIP(addr)
			if naddr.Equal(pubaddr) && auth.Allow(pubaddr) {
				a, err := admin.GetAdminFromFile(adminfile)
				if err == nil {
					s, err := a.Session()
					if err == nil {
						s.Addr = gen
						defer s.Close()
						info, err := s.AddTunnelIfNotThere(pubkey)
						if err == nil {
							info.Pubkey = pk
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
