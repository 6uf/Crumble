package src

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/6uf/apiGO"
	"github.com/6uf/h2"
	"github.com/PuerkitoBio/goquery"
	tls2 "github.com/refraction-networking/utls"
)

func ReturnPayload(acc, bearer, name string) string {
	if acc == "Giftcard" {
		return fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n{\"profileName\":\"%v\"}\r\n", len(`{"profileName":"`+name+`"}`), bearer, name)
	} else {
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: bearer " + bearer + "\r\n\r\n"
	}
}

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
}

func IsAvailable(name string) bool {
	resp, _ := http.Get("https://account.mojang.com/available/minecraft/" + name)
	return resp.StatusCode == 200
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

func ReturnOne(t1 time.Time, err error) TimeFluc { return TimeFluc{T1: t1, Err1: err} }
func OnlyAttr(val string, _ bool) string         { return val }
func OnlyTime(t time.Time, _ error) time.Time    { return t }
func ReturnAll(f T) (time.Time, error, time.Time, error) {
	return f.F[0].T1, f.F[0].Err1, f.F[1].T1, f.F[1].Err1
}
func Rewrite(path, accounts string) {
	os.Create(path)

	file, _ := os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()

	file.WriteAt([]byte(accounts), 0)
}

func GetDroptimes(name string) (int64, int64) {
	start := time.Now()
	var begin, end int64
	go func() {
		if conn, err := (&h2.Client{Config: h2.GetDefaultConfig()}).Connect("https://namemc.com/search?q="+name, h2.ReqConfig{ID: 1, BuildID: tls2.HelloChrome_100}); err == nil {
			if resp, err := conn.Do(h2.MethodGet, "", "", nil); err == nil && resp.Status == "200" {
				doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Data))
				if b, ok := doc.Find(`#availability-time`).Attr("datetime"); ok {
					if e, ok := doc.Find(`#availability-time2`).Attr("datetime"); ok {
						if t1, err := time.Parse(time.RFC3339, b); err == nil {
							if t2, err := time.Parse(time.RFC3339, e); err == nil {
								begin, end = t1.Unix(), t2.Unix()
							}
						}
					}
				}
			}
		}
	}()
	for {
		if begin != 0 && end != 0 {
			return begin, end
		} else {
			if time.Since(start) > 10*time.Second {
				return 0, 0
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func Update(new apiGO.Info) {
	for e, acc := range Bearer.Details {
		if acc.Email == new.Email && new.Error == "" {
			Bearer.Details[e] = new
			c := Con.Bearers[e]
			if c.Email == new.Email {
				Con.Bearers[e] = apiGO.Bearers{
					Bearer:       new.Bearer,
					Email:        new.Email,
					Password:     new.Password,
					AuthInterval: 86400,
					AuthedAt:     time.Now().Unix(),
					Type:         new.AccountType,
					NameChange:   true,
					Info:         new.Info,
				}
			}
			Con.SaveConfig()
			Con.LoadState()
		}
	}
}
