package admin

import (
	"testing"
)

func TestAdminPing(t *testing.T) {
	s, err := NewSession("127.0.0.1:11234")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer s.Close()
	err = s.Ping()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func sessionOrDie(t *testing.T) *Session {
	a, err := GetAdmin()
	if err != nil {
		t.Log(err)
		t.Fail()
		return nil
	}
	s, err := a.Session()
	if err != nil {
		t.Log(err)
		t.Fail()
		return nil
	}
	return s

}

func TestAdminAuthedPing(t *testing.T) {
	s := sessionOrDie(t)
	if s == nil {
		return
	}
	defer s.Close()
	ping := map[string]interface{}{"q": "ping"}
	r, err := s.Authed(ping)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if r["q"] != "pong" {
		t.Fail()
	}
}

func TestAddIpTunnel(t *testing.T) {
	s := sessionOrDie(t)
	if s == nil {
		return
	}
	key := "x68cm5b1tg7bdy6trlwgcg5gn9u7ttu6856g16k6dl963uz0v410.k"
	info := &IPTunnel{
		Pubkey:  key,
		Address: "10.1.1.1",
	}
	err := s.AddIPTunnel(info)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = s.RemoveIPTunnelsByPubkey(key)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log("removed okay")
	}
}

func TestListIpTunnels(t *testing.T) {

	s := sessionOrDie(t)
	if s == nil {
		return
	}
	defer s.Close()
	tuns, err := s.ListIPTunnels()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	for _, info := range tuns {
		t.Log(info)
	}
}
