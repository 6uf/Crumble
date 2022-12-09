package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/6uf/apiGO"
	"github.com/6uf/h2"
	"github.com/PuerkitoBio/goquery"
	html "github.com/antchfx/htmlquery"
	"github.com/iskaa02/qalam/gradient"
	tls2 "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
)

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
}

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
	if req, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", ip, port), &proxy.Auth{
		User:     user,
		Password: pass,
	}, proxy.Direct); err == nil {
		if conn, err := req.Dial("tcp", "api.minecraftservices.com:443"); err == nil {
			return SniperProxy{
				Proxy:        tls.Client(conn, &tls.Config{RootCAs: Roots, InsecureSkipVerify: true, ServerName: "api.minecraftservices.com"}),
				Alive:        true,
				ProxyDetails: Proxies{IP: ip, Port: port, User: user, Password: pass},
			}
		}
	}
	return SniperProxy{}
}

func SendWebhook(name, bearer string) {
	type Payload struct {
		Name   string `json:"name"`
		Bearer string `json:"bearer"`
		Url    string `json:"icon_url"`
	}
	http.Post(fmt.Sprintf("https://buxflip.com/data/webhook/%v/%v", Con.DiscordID, name), "application/json", bytes.NewBuffer(apiGO.JsonValue(Payload{Name: name, Bearer: bearer, Url: GetHeadUrl(name)})))
}

func GetHeadUrl(name string) string {
	if resp, err := http.Get("https://buxflip.com/data/namemc/head/" + name); err != nil {
		return "https://s.namemc.com/2d/skin/face.png?id=23ba96021149f38e&scale=4"
	} else if resp.StatusCode == 200 {
		var J struct {
			Head string `json:"headurl"`
		}
		json.Unmarshal([]byte(apiGO.ReturnJustString(io.ReadAll(resp.Body))), &J)
		if J.Head == "" {
			return "https://s.namemc.com/2d/skin/face.png?id=23ba96021149f38e&scale=4"
		}
		return J.Head
	}
	return "https://s.namemc.com/2d/skin/face.png?id=23ba96021149f38e&scale=4"
}

func IsAvailable(name string) bool {
	resp, err := http.Get("https://account.mojang.com/available/minecraft/" + name)
	if err == nil {
		return resp.StatusCode == 200
	} else {
		return false
	}
}

func GetDroptimes(name string) (int64, int64, string, string) {
	var st, se string
	if conn, err := (&h2.Client{Config: h2.GetDefaultConfig()}).Connect("https://namemc.com/search?q="+name, h2.ReqConfig{ID: 1, BuildID: tls2.HelloChrome_100, DataBodyMaxLength: 167859}); err == nil {
		if resp, err := conn.Do(h2.MethodGet, "", "", nil); err == nil && resp.Status == "200" {
			doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Data))
			if d, err := html.Parse(strings.NewReader(string(resp.Data))); err == nil {
				if status, searches := html.FindOne(d, `//*[@id="status-bar"]/div/div[1]/div[2]`), html.FindOne(d, `//*[@id="status-bar"]/div/div[2]/div[2]`); status.FirstChild != nil && searches.FirstChild != nil {
					st = status.FirstChild.Data
					se = strings.Split(searches.FirstChild.Data, " / month")[0]
				}
			}
			if b, ok := doc.Find(`#availability-time`).Attr("datetime"); ok {
				if e, ok := doc.Find(`#availability-time2`).Attr("datetime"); ok {
					if t1, err := time.Parse(time.RFC3339, b); err == nil {
						if t2, err := time.Parse(time.RFC3339, e); err == nil {

							return t1.Unix(), t2.Unix(), st, se
						}
					}
				}
			}
		}
	}
	return 0, 0, st, se
}

func WriteToLogs(name, logs string) {
	name = strings.ToLower(name)
	body, err := os.ReadFile("logs/names/" + name + ".txt")
	if os.IsNotExist(err) {
		os.Create("logs/names/" + name + ".txt")
	}
	str := string(body)
	str += logs
	os.WriteFile("logs/names/"+name+".txt", []byte(str), 0644)
}

func Logo(Data string) string {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors(RGB...).
		Mode(gradient.BlendRgb).
		Build()
	return g.Mutline(Data)
}
