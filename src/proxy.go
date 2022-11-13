package src

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
}

var Test []SniperProxy

func SendAlive() {
	for {
		for i, conns := range Test {
			go func(conns SniperProxy, i int) {
				if conns.Alive {
					fmt.Fprintf(conns.Proxy, "GET /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Abysal/1.1\r\n\r\n")
					buf := make([]byte, 4000)
					conns.Proxy.Read(buf)
					if string(buf)[9:12] == "" {
						Test[i].Alive = false
					}
				}
			}(conns, i)
		}
		time.Sleep(10 * time.Second)
	}
}

func GenProxysAndStoreForUsage() {
	var wg sync.WaitGroup
	for i, acc := range Proxy.Proxys {
		wg.Add(1)
		go func(acc string, i int) {
			if con := Connect(acc); con.Proxy != nil {
				Test = append(Test, con)
			}
			wg.Done()
		}(acc, i)
	}
	wg.Wait()
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

func RandomProxyConn() SniperProxy {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	P := Test[rand.Intn(len(Test))]
	if P.Alive {
		return P
	}
	return SniperProxy{}
}
