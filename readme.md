# cjdns exit tunnel broker


## server setup

build server and run it at http://[your:mesh:net:ip]:port/

    $ go install github.com/majestrate/tuntun/cmd/tuntun-server
    $ $GOPATH/bin/tuntun-server port
    
## client setup

build client and obtain a nat address from a tuntun server already set up

    $ go install github.com/majestrate/tuntun/cmd/tuntun-client
    $ $GOPATH/bin/tuntun-client http://[some:mesh:net:ip]:port/




