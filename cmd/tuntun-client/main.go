package main

import (
	"encoding/json"
	//	"github.com/majestrate/tuntun/lib/api"
	"github.com/majestrate/tuntun/lib/api/admin"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s http://someserver:port/", os.Args[0])
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	a, err := admin.GetAdmin()
	if err != nil {
		log.Fatal(err)
	}
	s, err := a.Session()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	log.Println("Trying to get address...")
	r, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == 200 {
		var info admin.IPTunnel
		err = json.NewDecoder(r.Body).Decode(&info)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Adding tunnel via", info.Pubkey, "using", info.Address)
		err = s.ConnectIPTunnel(info.Pubkey)
		if err == nil {
			log.Println("success")
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("HTTP %d", r.StatusCode)
	}
}
