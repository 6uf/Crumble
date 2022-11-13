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

func GenProxysAndStoreForUsage() {
	for {
		var wg sync.WaitGroup
		for i, acc := range Proxy.Proxys {
			wg.Add(1)
			go func(acc string, i int) {
				var Found bool
				if con := Connect(acc); con.Proxy != nil {
					for i, Con := range Test {
						if Con.ProxyDetails.IP == strings.Split(acc, ":")[0] {
							Found = true
							Test[i] = con
						}
					}
					if !Found {
						Test = append(Test, con)
					}
				}
				wg.Done()
			}(acc, i)
		}
		wg.Wait()
		time.Sleep(15 * time.Second)
	}
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
