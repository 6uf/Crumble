package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/utils"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/6uf/StrCmd"
	"github.com/6uf/apiGO"
)

var Username = ""

func TempCalc() time.Duration {
	var Start float64 = 0.1
	for {
		if Num := float64(len(utils.Bearer.Details)) * Start; Num >= 13000 {
			return time.Duration(Start) * time.Millisecond
		}
		Start++
	}
}

func init() {
	utils.Roots.AppendCertsFromPEM(utils.ProxyByte)
	apiGO.Clear()
	utils.Con.LoadState()
	fmt.Print(utils.Logo(`_________                        ______ ______     
__  ____/__________  ________ ______  /____  /____ 
_  /    __  ___/  / / /_  __ '__ \_  __ \_  /_  _ \
/ /___  _  /   / /_/ /_  / / / / /  /_/ /  / /  __/
\____/  /_/    \__,_/ /_/ /_/ /_//_.___//_/  \___/ 

`))

	if utils.Con.DiscordID == "" {
		fmt.Print(utils.Logo("Discord ID: "))
		fmt.Scan(&utils.Con.DiscordID)
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}
	if utils.Con.FirstUse {
		fmt.Print(utils.Logo("\nUse proxys for authentication? : [YES/NO] > "))
		var ProxyAuth string
		fmt.Scan(&ProxyAuth)
		utils.Con.FirstUse = false
		utils.Con.UseProxyDuringAuth = strings.Contains(strings.ToLower(ProxyAuth), "y")
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}
	if utils.Con.FirstUse {
		fmt.Print(utils.Logo("\nSent to webhook when a snipe occurs? : [YES/NO] > "))
		var WebhookCheck string
		fmt.Scan(&WebhookCheck)
		utils.Con.FirstUse = false
		utils.Con.SendWebhook = strings.Contains(strings.ToLower(WebhookCheck), "y")
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}
	Username = GetDiscordUsername(utils.Con.DiscordID)
	if file_name := "accounts.txt"; utils.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if file_name := "names.txt"; utils.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if file_name := "proxys.txt"; utils.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.MkdirAll("logs/names", 0755)
	}
	utils.Proxy.GetProxys(false, nil)
	utils.Proxy.Setup()
	utils.AuthAccs()
	go utils.CheckAccs()
	fmt.Println(utils.Logo(fmt.Sprintf("<!> accs available    [%v]\n<!> proxies available [%v]\n", len(utils.Bearer.Details), len(utils.Proxy.Proxys))))
}

func main() {
	app := StrCmd.App{
		Version:        "v1.0.25b-CR",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					ReqAmt := 0
					Claimed := false
					name := StrCmd.String("-u")
					Spread := TempCalc()

					start, end := utils.GetDroptimes(name)
					drop := time.Unix(start, 0)
					for time.Now().Before(drop) {
						fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v", name, time.Until(drop).Round(time.Second)))))
						time.Sleep(time.Second * 1)
					}

					c := make(chan os.Signal, 1)
					signal.Notify(c, os.Interrupt)
					go func() {
						<-c
						signal.Stop(c)
						Claimed = true
					}()
					go func() {
						for {
							if utils.IsAvailable(name) {
								Claimed = true
								break
							}
							if time.Now().After(time.Unix(end, 0)) {
								Claimed = true
								break
							}
							time.Sleep(10 * time.Second)
						}
					}()
					if name != "" {
						if len(utils.Bearer.Details) != 0 {
							for _, Config := range utils.Bearer.Details {
								go func(Config apiGO.Info, name string) {
									Next := time.Now()
									for !Claimed {
										New := time.Now().Add(time.Second * 15)
										if time.Until(Next).Seconds() > 10 {
											time.Sleep(10 * time.Second)
										}
										for _, Acc := range utils.Bearer.Details {
											if strings.EqualFold(Acc.Email, Config.Email) {
												Config = Acc
												break
											}
										}
										if proxy := utils.Connect(utils.Proxy.CompRand()); proxy.Alive {
											Payload := ReturnPayload(Config.AccountType, Config.Bearer, name)
											fmt.Fprint(proxy.Proxy, Payload[:len(Payload)-4])
											time.Sleep(time.Until(Next))
											Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(proxy.Proxy, "\r\n"), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType, Info: Config.Info}
											ReqAmt++
											C := fmt.Sprintf("[%v >> %v] (%v:%v) %v - %v (%v)", Req.ResponseDetails.SentAt.Format("15:04:05.0000"), Req.ResponseDetails.RecvAt.Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, name, Req.Email[:4], proxy.ProxyDetails.IP[:5]+strings.Repeat("*", 5), ReqAmt)
											fmt.Println(utils.Logo(C))
											utils.WriteToLogs(name, C+"\n")
											switch Req.ResponseDetails.StatusCode {
											case "200":
												Claimed = true
												if utils.Con.SkinChange.Link != "" {
													apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Req.Bearer)
												}
												if utils.Con.SendWebhook {
													go utils.SendWebhook(name, Req.Bearer)
												}
												fmt.Printf("[%v] Succesful - %v %v\n", name, Req.Email, apiGO.NameMC(Req.Bearer, apiGO.GetProfileInformation(Req.Bearer, (*apiGO.ProxyMS)(&proxy.ProxyDetails))))
											}
										}
										Next = New
									}
								}(Config, name)
								time.Sleep(Spread)
							}
							for !Claimed {
								time.Sleep(10 * time.Second)
							}
							signal.Stop(c)
						} else {
							fmt.Println(utils.Logo(fmt.Sprintf("Unable to start process for %v, no bearers found.", StrCmd.String("-u"))))
						}
					}
				},
				Args: []string{
					"-u",
				},
			},
			"list": {
				Description: "List snipes from accounts within the names.txt file and send req at random based on each.",
				Action: func() {
					// Get accounts
					Spread := TempCalc()
					accs, _ := os.ReadFile("names.txt")
					Scanner := bufio.NewScanner(bytes.NewBuffer(accs))
					type Names struct {
						Name  string
						Taken bool
					}
					var Accs []Names
					for Scanner.Scan() {
						if Text := Scanner.Text(); Text != "" {
							Accs = append(Accs, Names{
								Name: Text,
							})
						}
					}
					ReqAmt := 0
					c := make(chan os.Signal, 1)
					signal.Notify(c, os.Interrupt)
					go func() {
						<-c
						signal.Stop(c)
						for i := range Accs {
							Accs[i].Taken = true
						}
					}()
					go func() {
						for {
							for i, n := range Accs {
								if !n.Taken && utils.IsAvailable(n.Name) {
									Accs[i].Taken = true
								}
								time.Sleep(10 * time.Second)
							}
						}
					}()
					if len(utils.Bearer.Details) != 0 {
						for _, Config := range utils.Bearer.Details {
							go func(Config apiGO.Info) {
								Next := time.Now()
								for {
									New := time.Now().Add(time.Second * 15)
									if time.Until(Next).Seconds() > 10 {
										time.Sleep(10 * time.Second)
									}
									for _, Acc := range utils.Bearer.Details {
										if strings.EqualFold(Acc.Email, Config.Email) {
											Config = Acc
											break
										}
									}
									if proxy := utils.Connect(utils.Proxy.CompRand()); proxy.Alive {
										rand.Seed(time.Now().UnixMilli())
										if Data := Accs[rand.Intn(len(Accs))]; !Data.Taken {
											name := Data.Name
											Payload := ReturnPayload(Config.AccountType, Config.Bearer, name)
											fmt.Fprint(proxy.Proxy, Payload[:len(Payload)-4])
											time.Sleep(time.Until(Next))
											Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(proxy.Proxy, "\r\n"), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType, Info: Config.Info}
											ReqAmt++
											C := fmt.Sprintf("[%v >> %v] (%v:%v) %v - %v (%v)", Req.ResponseDetails.SentAt.Format("15:04:05.0000"), Req.ResponseDetails.RecvAt.Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, name, Req.Email[:4], proxy.ProxyDetails.IP[:5]+strings.Repeat("*", 5), ReqAmt)
											fmt.Println(utils.Logo(C))
											utils.WriteToLogs(name, C+"\n")
											switch Req.ResponseDetails.StatusCode {
											case "200":
												for i, n := range Accs {
													if strings.EqualFold(n.Name, name) {
														Accs[i].Taken = true
														break
													}
												}
												if utils.Con.SkinChange.Link != "" {
													apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Req.Bearer)
												}
												if utils.Con.SendWebhook {
													go utils.SendWebhook(name, Req.Bearer)
												}
												fmt.Printf("[%v] Succesful - %v %v\n", name, Req.Email, apiGO.NameMC(Req.Bearer, apiGO.GetProfileInformation(Req.Bearer, (*apiGO.ProxyMS)(&proxy.ProxyDetails))))
											}
										} else {
											Found := 0
											for _, n := range Accs {
												if !n.Taken {
													Found++
												}
											}
											if Found == 0 {
												break
											}
										}
									}
									Next = New
								}
							}(Config)
							time.Sleep(Spread)
						}
						for {
							Found := 0
							for _, n := range Accs {
								if !n.Taken {
									Found++
								}
							}
							if Found == 0 {
								break
							}
							time.Sleep(10 * time.Second)
						}
					} else {
						fmt.Println(utils.Logo(fmt.Sprintf("Unable to start process for %v, no bearers found.", StrCmd.String("-u"))))
					}
				},
			},
			"key": {
				Description: "Gets your namemc key!",
				Action: func() {
					var details string
					fmt.Print(utils.Logo("[email:pass] > "))
					fmt.Scan(&details)
					if acc := strings.Split(details, ":"); len(acc) > 0 {
						if len(utils.Proxy.Proxys) > 0 {
							ip, port, user, pass := "", "", "", ""
							switch data := strings.Split(utils.Proxy.CompRand(), ":"); len(data) {
							case 2:
								ip = data[0]
								port = data[1]
							case 4:
								ip = data[0]
								port = data[1]
								user = data[2]
								pass = data[3]
							}
							Acc := apiGO.MS_authentication(acc[0], acc[1], &apiGO.ProxyMS{
								IP: ip, Port: port, User: user, Password: pass,
							})
							fmt.Println(utils.Logo(apiGO.NameMC(Acc.Bearer, Acc.Info)))
						} else {
							Acc := apiGO.MS_authentication(acc[0], acc[1], nil)
							fmt.Println(utils.Logo(apiGO.NameMC(Acc.Bearer, Acc.Info)))
						}
					}
				},
			},
		},
	}
	app.Run(utils.Logo(fmt.Sprintf("@%v/root: ", Username)))
}

func ReturnPayload(acc, bearer, name string) string {
	if acc == "Giftcard" {
		return fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n{\"profileName\":\"%v\"}\r\n", len(`{"profileName":"`+name+`"}`), bearer, name)
	} else {
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: bearer " + bearer + "\r\n\r\n"
	}
}

func GetDiscordUsername(ID string) string {
	resp, err := http.Get("https://buxflip.com/data/discord/" + ID)
	if err != nil {
		return "Unknown"
	} else {
		if resp.StatusCode == 429 {
			return "Unknown"
		}
		var Body struct {
			Data struct {
				Name string `json:"username"`
			} `json:"data"`
		}
		json.Unmarshal([]byte(apiGO.ReturnJustString(io.ReadAll(resp.Body))), &Body)
		return Body.Data.Name
	}
}
