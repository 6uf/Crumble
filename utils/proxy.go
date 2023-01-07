package utils

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
)

func Connect(proxy string) (*tls.Conn, bool) {
	ip := strings.Split(proxy, ":")
	if conn, err := net.Dial("tcp", ip[0]+":"+ip[1]); err == nil {
		if len(ip) > 2 {
			conn.Write([]byte(fmt.Sprintf("CONNECT minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net:443 HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net:443\r\nProxy-Authorization: Basic %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\n\r\n", base64.RawStdEncoding.EncodeToString([]byte(ip[2]+":"+ip[3])))))
		} else {
			conn.Write([]byte("CONNECT minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net:443 HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net:443\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\n\r\n"))
		}
		var junk = make([]byte, 4096)
		conn.Read(junk)
		switch Status := string(junk); Status[9:12] {
		case "200":
			return tls.Client(conn, &tls.Config{RootCAs: Roots, InsecureSkipVerify: true, ServerName: "minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net"}), true
		case "407":
			fmt.Println(Logo(fmt.Sprintf("[%v] Proxy <%v> Failed to authorize: Username/Password invalid.", Status, ip[0])))
		default:
			return nil, false
		}
	}
	return nil, false
}
