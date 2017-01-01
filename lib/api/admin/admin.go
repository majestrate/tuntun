package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Admin struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func GetAdmin() (a *Admin, err error) {
	h := os.Getenv("HOME")
	fpath := filepath.Join(h, ".cjdnsadmin")
	return GetAdminFromFile(fpath)
}

func GetAdminFromFile(path string) (a *Admin, err error) {
	var d []byte
	d, err = ioutil.ReadFile(path)
	if err == nil {
		a = new(Admin)
		err = json.Unmarshal(d, a)
	}
	return
}

func (a *Admin) Session() (*Session, error) {
	s, err := NewSession(fmt.Sprintf("%s:%d", a.Addr, a.Port))
	if err == nil {
		s.p = a.Password
	}
	return s, err
}
