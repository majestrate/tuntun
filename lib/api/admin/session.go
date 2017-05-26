package admin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/majestrate/tuntun/lib/util"
	"github.com/zeebo/bencode"
	"io"
	"net"
)

var ErrBadAuth = errors.New("bad auth")
var ErrBadResp = errors.New("bad response")
var ErrBadRepl = errors.New("bad reply")
var ErrBadAddr = errors.New("bad admin address")
var ErrDuplicateAddress = errors.New("duplicate address")

type Session struct {
	a    net.Addr
	c    net.PacketConn
	p    string
	Addr util.AddressGenerator
}

func NewSession(addr string) (*Session, error) {
	u, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	lu, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", lu)
	if err != nil {
		return nil, err
	}
	return &Session{
		a: u,
		c: conn,
	}, nil
}

func (s *Session) Close() (err error) {
	err = s.c.Close()
	return
}

func (s *Session) sendRaw(obj map[string]interface{}) (err error) {
	var d []byte
	d, err = bencode.EncodeBytes(obj)
	if err == nil {
		_, err = s.c.WriteTo(d, s.a)
	}
	return
}

func (s *Session) recvRaw(obj *map[string]interface{}) (err error) {
	var d [65536]byte
	var n int
	n, _, err = s.c.ReadFrom(d[:])
	if err == nil {
		err = bencode.DecodeBytes(d[:n], obj)
	}
	return
}

// build authenticated request
func (s *Session) buildReq(passwd string, authreq map[string]interface{}) (built map[string]interface{}, err error) {
	// get cookie
	var cookie_resp map[string]interface{}
	cookie_resp, err = s.Command(map[string]interface{}{"q": "cookie"})
	if err == nil {
		c, ok := cookie_resp["cookie"]
		if ok {
			cookie := c.(string)
			d := sha256.Sum256([]byte(passwd + cookie))
			built = map[string]interface{}{
				"q":      "auth",
				"hash":   hex.EncodeToString(d[:]),
				"cookie": cookie,
				"aq":     authreq["q"],
			}
			for k, v := range authreq {
				if k == "q" || k == "cookie" || k == "aq" || k == "hash" {
					continue
				}
				built[k] = v
			}
			var data []byte
			data, err = bencode.EncodeBytes(built)
			if err == nil {
				fd := sha256.Sum256(data)
				built["hash"] = hex.EncodeToString(fd[:])
			}
		} else {
			err = ErrBadRepl
		}
	}
	return
}

// build a transaction
func (s *Session) buildTx(req map[string]interface{}) (built map[string]interface{}) {
	req["txid"] = util.RandStr(10)
	built = req
	return
}

// send a command and get the reply
func (s *Session) Command(obj map[string]interface{}) (response map[string]interface{}, err error) {
	resp := make(map[string]interface{})
	err = s.sendRaw(obj)
	if err == nil {
		err = s.recvRaw(&resp)
		if err == nil {
			e, ok := resp["error"]
			if ok && e != "none" {
				err = errors.New(e.(string))
			} else {
				response = resp
			}
		}
	}
	return
}

// send a ping
func (s *Session) Ping() (err error) {
	ping := map[string]interface{}{"q": "ping"}
	var resp map[string]interface{}
	resp, err = s.Command(ping)
	if err == nil {
		pong, ok := resp["q"]
		if !ok || pong != "pong" {
			err = ErrBadResp
		}
	}
	return
}

// do an authenticated command
func (s *Session) Authed(obj map[string]interface{}) (response map[string]interface{}, err error) {
	obj, err = s.buildReq(s.p, obj)
	if err == nil {
		response, err = s.Command(obj)
	}
	return
}

func (s *Session) ConnectIPTunnel(pubkey string) (err error) {
	_, err = s.Authed(map[string]interface{}{
		"q": "IpTunnel_connectTo",
		"args": map[string]interface{}{
			"publicKeyOfNodeToConnectTo": pubkey,
		},
	})
	return
}

func (s *Session) AddTunnelIfNotThere(pubkey string) (info *IPTunnel, err error) {
	var infos []*IPTunnel
	addrs := make(map[string]*IPTunnel)
	infos, err = s.ListIPTunnels()
	if err == nil {
		for idx := range infos {
			if infos[idx].Pubkey == pubkey {
				info = infos[idx]
			}
			a := infos[idx].Address
			if a != "" {
				addrs[a] = infos[idx]
			}
		}
		if info == nil {
			if s.Addr == nil {
				for {
					// make unique address
					var n [3]byte
					io.ReadFull(rand.Reader, n[:])
					a := net.IPv4(10, n[0], n[1], n[2]).String()
					_, ok := addrs[a]
					if !ok {
						info = &IPTunnel{
							Address: a,
							Pubkey:  pubkey,
						}
						err = s.AddIPTunnel(info)
						break
					}
				}
			} else {
				var ip net.IP
				ip, err = s.Addr.GetIPFor(pubkey)
				if err == nil {
					a := ip.String()
					_, ok := addrs[a]
					if ok {
						err = ErrDuplicateAddress
					} else {
						info = &IPTunnel{
							Address: a,
							Pubkey:  pubkey,
						}
					}
				}
			}
		}
	}
	return
}

func (s *Session) AddIPTunnel(info *IPTunnel) (err error) {
	if info == nil {
		err = errors.New("info was nil")
		return
	}
	_, err = s.Authed(map[string]interface{}{
		"q": "IpTunnel_allowConnection",
		"args": map[string]interface{}{
			"ip4Prefix":                 32,
			"ip4Alloc":                  32,
			"ip4Address":                info.Address,
			"publicKeyOfAuthorizedNode": info.Pubkey,
		},
	})
	return
}

func (s *Session) GetIPTunnel(idx int64) (info *IPTunnel, err error) {
	var r map[string]interface{}
	r, err = s.Authed(map[string]interface{}{
		"q": "IpTunnel_showConnection",
		"args": map[string]interface{}{
			"connection": idx,
		},
	})
	if err == nil {
		addr, ok := r["ip4Address"]
		if !ok {
			err = errors.New("no ip4 addresss")
			return
		}
		key, ok := r["key"]
		if !ok {
			err = errors.New("no pubkey")
			return
		}
		info = &IPTunnel{
			Address: addr.(string),
			Pubkey:  key.(string),
		}
	}
	return
}

func (s *Session) RemoveIPTunnel(idx int64) (err error) {
	_, err = s.Authed(map[string]interface{}{
		"q": "IpTunnel_removeConnection",
		"args": map[string]interface{}{
			"connection": idx,
		},
	})
	return
}

func (s *Session) RemoveIPTunnelsByPubkey(key string) (err error) {
	var r map[string]interface{}
	r, err = s.Authed(map[string]interface{}{
		"q": "IpTunnel_listConnections",
	})
	if err == nil {
		var conns []interface{}
		c, ok := r["connections"]
		if !ok {
			err = errors.New("no connections in response")
			return
		}
		conns = c.([]interface{})
		for _, idx := range conns {
			var info *IPTunnel
			info, err = s.GetIPTunnel(idx.(int64))
			if err == nil && info.Pubkey == key {
				err = s.RemoveIPTunnel(idx.(int64))
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (s *Session) ListIPTunnels() (tunnels []*IPTunnel, err error) {
	var r map[string]interface{}
	r, err = s.Authed(map[string]interface{}{
		"q": "IpTunnel_listConnections",
	})
	if err == nil {
		var conns []interface{}
		c, ok := r["connections"]
		if !ok {
			err = errors.New("no connections in response")
			return
		}
		conns = c.([]interface{})
		for _, idx := range conns {
			var info *IPTunnel
			info, err = s.GetIPTunnel(idx.(int64))
			if err == nil {
				tunnels = append(tunnels, info)
			} else {
				return
			}
		}
	}
	return
}

func (s *Session) GetOurPubkey() (key string, err error) {
	var r map[string]interface{}
	r, err = s.Authed(map[string]interface{}{
		"q": "NodeStore_nodeForAddr",
		"args": map[string]interface{}{
			"ip": 0,
		},
	})
	if err == nil {
		res := r["result"].(map[string]interface{})
		key = res["key"].(string)
	}
	return
}
