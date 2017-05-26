package config

import (
	"github.com/majestrate/configparser"
	"github.com/majestrate/tuntun/lib/util"
	"net"
	"strconv"
)

type Config struct {
	Port int
	Exit ExitConfig
	Auth AuthConfig
}

type ExitConfig struct {
	Strategy string
	Net      *net.IPNet
}

func (c *ExitConfig) AddressGenerator() util.AddressGenerator {
	return nil
}

type AuthConfig struct {
}

func (c *AuthConfig) AuthPolicy() util.ClientAuthPolicy {
	return &util.NullAuthPolicy{}
}

func Load(fname string) (conf *Config, err error) {
	var c *configparser.Configuration
	c, err = configparser.Read(fname)
	if err == nil {
		var s *configparser.Section
		s, err = c.Section("server")
		if err != nil {
			return
		}
		port := s.ValueOf("port")
		if port == "" {
			port = "1800"
		}
		var p int
		p, err = strconv.Atoi(port)
		if err != nil {
			return
		}
		conf = &Config{
			Port: p,
		}
		_, conf.Exit.Net, err = net.ParseCIDR(s.ValueOf("exit-net"))
		if err != nil {
			conf = nil
			return
		}
	}
	return
}
