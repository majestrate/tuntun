# cjdns exit tunnel broker


## server setup

build server and run it at http://[your:mesh:net:ip]:port/

    $ go get -u github.com/majestrate/tuntun/cmd/tuntun-server
    $ $GOPATH/bin/tuntun-server port
    
uses `10.0.0.0/8` for addresses
    
    # echo 1 > /proc/sys/net/ipv4/ip_forward
    # iptables -t nat -A POSTROUTING -s 10.0.0.0/8 -o eth0 -j MASQUERADE
    # ip route add 10.0.0.0/8 dev tun0


## client setup

build client and obtain a nat address from a tuntun server already set up

    $ go get -u github.com/majestrate/tuntun/cmd/tuntun-client
    $ $GOPATH/bin/tuntun-client http://[some:mesh:net:ip]:port/
    
push your route over cjdns

    # ip route add default dev tun0


## TODO

* make server address range configurable
* add client authentication
