package src

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
}

func Connect(acc string) SniperProxy {
	var user, pass, ip, port string
	auth := strings.Split(acc, ":")
	switch len(auth) {
	case 2:
		ip, port = auth[0], auth[1]
	case 4:
		ip, port, user, pass = auth[0], auth[1], auth[2], auth[3]
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(ProxyByte)
	req, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", ip, port), &proxy.Auth{
		User:     user,
		Password: pass,
	}, proxy.Direct)
	if err == nil {
		if conn, err := req.Dial("tcp", "api.minecraftservices.com:443"); err == nil {
			cd, _ := conn.(*net.TCPConn)
			cd.SetKeepAlive(true)
			return SniperProxy{
				Proxy:        tls.Client(cd, &tls.Config{RootCAs: roots, InsecureSkipVerify: true, ServerName: "api.minecraftservices.com"}),
				UsedAt:       time.Now(),
				Alive:        true,
				ProxyDetails: Proxies{IP: ip, Port: port, User: user, Password: pass},
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return SniperProxy{}
}
