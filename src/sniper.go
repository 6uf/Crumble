package src

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/6uf/apiGO"
)

func SnipeDefault(name string) {
	var Accs [][]apiGO.Info
	var incr int
	var use int
	for _, Acc := range Bearer.Details {
		if len(Accs) == 0 {
			Accs = append(Accs, []apiGO.Info{
				Acc,
			})
		} else {
			if incr == 3 {
				incr = 0
				use++
				Accs = append(Accs, []apiGO.Info{})
			}
			Accs[use] = append(Accs[use], Acc)
		}
		incr++
	}

	start, _ := GetDroptimes(name)
	GotName, Terminate, Taken, c := make(chan string), false, IsAvailable(name), make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ClaimedIGN := false

	go func() {
		<-c
		signal.Stop(c)
		Terminate = true
	}()

	drop := time.Unix(start, 0)
	for time.Now().Before(drop) {
		fmt.Printf("[%v] %v    \r", name, time.Until(drop).Round(time.Second))
		time.Sleep(time.Second * 1)
	}

	go func() {
		if !Taken {
			for !IsAvailable(name) {
				time.Sleep(10 * time.Second)
			}
			if !ClaimedIGN {
				GotName <- fmt.Sprintf("[%v] has become unavailable [%v]", name, time.Now().Unix())
				Taken = true
				signal.Stop(c)
			}
		} else {
			GotName <- fmt.Sprintf("[%v] has become unavailable [%v]", name, time.Now().Unix())
		}
	}()

	if len(Proxy.Proxys) != 0 {
		go func() {
			ReqAmt := 0
			for !Terminate {
				if ClaimedIGN {
					return
				}
				for e, acc := range Accs {
					go func(acc []apiGO.Info) {
						for i, Config := range acc {
							if !Taken {
								if proxy := Connect(Proxy.CompRand()); proxy.Alive {
									for o := 0; int64(o) < Con.ReqAmtPerAcc; o++ {
										Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(proxy.Proxy, ReturnPayload(Config.AccountType, Config.Bearer, name)), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType, Info: Config.Info}
										ReqAmt++
										body := fmt.Sprintf("[%v] (%v:%v) %v - %v (%v)\n", time.Now().Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, name, Req.Email[:4], proxy.ProxyDetails.IP[:5]+strings.Repeat("*", 5), ReqAmt)
										fmt.Print(body)
										WriteToLogs(name, body)
										switch Req.ResponseDetails.StatusCode {
										case "200":
											ClaimedIGN = true
											if Con.SkinChange.Link != "" {
												apiGO.ChangeSkin(apiGO.JsonValue(Con.SkinChange), Req.Bearer)
											}
											fmt.Printf("[%v] Succesful - %v %v\n", name, Req.Email, apiGO.NameMC(Req.Bearer))
											SendWebhook(name, Req.Bearer)
											GotName <- fmt.Sprintf("[%v] Succesful - %v %v", name, Req.Email, apiGO.NameMC(Req.Bearer))
											signal.Stop(c)
											new, list, Accz := []apiGO.Bearers{}, []apiGO.Info{}, []string{}
											for _, Accs := range Con.Bearers {
												if Req.Email != Accs.Email {
													new = append(new, Accs)
												}
											}
											Con.Bearers = new
											for _, Accs := range Con.Bearers {
												for _, Acc := range Bearer.Details {
													if Acc.Email != Accs.Email {
														list = append(list, Acc)
													}
												}
											}
											Bearer.Details = list
											file, _ := os.Open("accounts.txt")
											defer file.Close()
											scanner := bufio.NewScanner(file)
											for scanner.Scan() {
												if strings.Split(scanner.Text(), ":")[0] != Req.Email {
													Accz = append(Accz, scanner.Text())
												}
											}
											Rewrite("accounts.txt", strings.Join(Accz, "\n"))
											Con.SaveConfig()
											Con.LoadState()
										case "401":
											fmt.Printf("[%v] %v came up invalid, reauthing..\n", Req.ResponseDetails.StatusCode, HashMessage(Req.Email, len(Req.Email)/4))
											Accs[e][i].Error = "AuthRequired:" + Config.Email
										}
										time.Sleep(time.Duration(Con.SpreadPerSend) * time.Millisecond)
									}
								}
							}
						}
					}(acc)
					for i, acc := range acc {
						if data := strings.Split(acc.Error, ":"); len(data) > 0 && data[0] == "AuthRequired" && data[1] == acc.Email {
							p := Proxy.CompRand()
							s := strings.Split(p, ":")
							var ip, port, user, pass string
							if len(s) > 2 {
								ip, port, user, pass = s[0], s[1], s[2], s[3]
							} else {
								ip, port = s[0], s[1]
							}
							upd := apiGO.MS_authentication(acc.Email, acc.Password, &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass})
							Update(upd)
							Accs[e][i] = upd
							break
						}
					}
					if Taken {
						GotName <- fmt.Sprintf("[%v] has become unavailable [%v]", name, time.Now().Unix())
						return
					}
					if Terminate {
						GotName <- "Terminated out of process for " + name
						return
					}
					if ClaimedIGN {
						return
					}
					time.Sleep(time.Duration(Con.SpreadPerAccount) * time.Millisecond)
				}
			}
		}()
	} else {
		fmt.Println("Cannot find any proxies, this sniper is only functional with proxys, as this is the private build.")
		GotName <- "Terminated out of process for " + name
	}
	fmt.Println(<-GotName)
}

func SendWebhook(name, bearer string) {
	type Payload struct {
		Name   string `json:"name"`
		Bearer string `json:"bearer"`
		Url    string `json:"icon_url"`
	}
	resp, _ := http.Post(fmt.Sprintf("https://buxflip.com/data/webhook/%v/%v", Con.DiscordID, name), "application/json", bytes.NewBuffer(apiGO.JsonValue(Payload{Name: name, Bearer: bearer, Url: GetHeadUrl(name)})))
	if resp.StatusCode == 200 {
		fmt.Println("Succesfully sent webhook!")
	} else {
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("Error: While sending webhook returned this body > " + string(body))
		}
	}
}

func GetHeadUrl(name string) string {
	if resp, err := http.Get("https://buxflip.com/data/namemc/head/" + name); err != nil {
		return "https://s.namemc.com/2d/skin/face.png?id=23ba96021149f38e&scale=4"
	} else if resp.StatusCode == 200 {
		var J struct {
			Head string `json:"headurl"`
		}
		json.Unmarshal([]byte(apiGO.ReturnJustString(io.ReadAll(resp.Body))), &J)
		return J.Head
	}
	return "https://s.namemc.com/2d/skin/face.png?id=23ba96021149f38e&scale=4"
}
