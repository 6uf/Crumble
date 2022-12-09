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
	return time.Duration(15000/len(utils.Bearer.Details)) * time.Millisecond
}

func init() {
	utils.Roots.AppendCertsFromPEM(utils.ProxyByte)
	apiGO.Clear()
	utils.Con.LoadState()
	for _, rgb := range utils.Con.Gradient {
		utils.RGB = append(utils.RGB, fmt.Sprintf("rgb(%v,%v,%v)", rgb.R, rgb.G, rgb.B))
	}
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
}

func main() {
	app := StrCmd.App{
		Version:        "v1.2.25b-CR",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					cl, name, Spread, Changed, Dummy, c, ChangeDetected := false, StrCmd.String("-u"), time.Millisecond, false, make(chan string), make(chan os.Signal, 1), make(chan apiGO.Details)
					if utils.Con.UseCustomSpread {
						Spread = time.Duration(utils.Con.Spread) * time.Millisecond
					} else {
						Spread = TempCalc()
					}
					signal.Notify(c, os.Interrupt)
					start, end, status, searches := utils.GetDroptimes(name)
					drop := time.Unix(start, 0)
					for time.Now().Before(drop) {
						fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v                 \r", name, time.Until(drop).Round(time.Second)))))
						time.Sleep(time.Second * 1)
					}
					go func() {
					Exit:
						for {
							select {
							case <-c:
								signal.Stop(c)
								Changed = true
								cl = true
								break Exit
							default:
								if utils.IsAvailable(name) {
									Changed = true
									cl = true
									break Exit
								}
								if time.Now().After(time.Unix(end, 0)) {
									Changed = true
									cl = true
									break Exit
								}
								time.Sleep(10 * time.Second)
							}
						}
					}()
					fmt.Print(utils.Logo(fmt.Sprintf(`
Name(s)    ~ %v
Spread     ~ %v
Proxies    ~ %v
Account(s) ~ %v
Searches   ~ %v
Status     ~ %v
Start      ~ %v
End        ~ %v

`, name, Spread, len(utils.Proxy.Proxys), len(utils.Bearer.Details), searches, status, time.Unix(start, 0), time.Unix(end, 0))))
					for _, Config := range utils.Bearer.Details {
						go Snipe(Config, name, &Changed, &ChangeDetected, false, nil, &Dummy)
						time.Sleep(Spread)
					}
				Exit:
					for {
						if cl {
							fmt.Println("\n", utils.Logo(name+" Has dropped."))
							signal.Stop(c)
							break Exit
						}
						select {
						case Req := <-ChangeDetected:
							if utils.Con.SkinChange.Link != "" {
								go apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Req.Bearer)
							}
							if utils.Con.SendWebhook {
								go utils.SendWebhook(name, Req.Bearer)
							}
							fmt.Println(utils.Logo(fmt.Sprintf("[%v] Succesfully sniped! - %v", name, Req.Email)))
							signal.Stop(c)
							break Exit
						default:
							time.Sleep(10 * time.Second)
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
					accs, _ := os.ReadFile("names.txt")
					var plain_accs []string
					Spread, ListName, Accs, Scanner, c, ChangeDetected := TempCalc(), make(chan string), []utils.Names{}, bufio.NewScanner(bytes.NewBuffer(accs)), make(chan os.Signal, 1), make(chan apiGO.Details)
					signal.Notify(c, os.Interrupt)
					for Scanner.Scan() {
						if Text := Scanner.Text(); Text != "" {
							plain_accs = append(plain_accs, Text)
							Accs = append(Accs, utils.Names{
								Name: Text,
							})
						}
					}
					go func() {
					Exit:
						for {
							select {
							case <-c:
								signal.Stop(c)
								for i := range Accs {
									ListName <- Accs[i].Name
								}
								break Exit
							default:
								for i, n := range Accs {
									if !n.Taken && utils.IsAvailable(n.Name) {
										ListName <- Accs[i].Name
									}
									time.Sleep(10 * time.Second)
								}
							}
						}
					}()
					fmt.Print(utils.Logo(fmt.Sprintf(`
Name(s)    ~ %v
Spread     ~ %v
Proxies    ~ %v
Account(s) ~ %v

`, plain_accs, Spread, len(utils.Proxy.Proxys), len(utils.Bearer.Details))))
					for _, Config := range utils.Bearer.Details {
						go Snipe(Config, "", &[]bool{false}[0], &ChangeDetected, true, &Accs, &ListName)
						time.Sleep(Spread)
					}
				Exit:
					for {
						select {
						case Req := <-ChangeDetected:
							if utils.Con.SkinChange.Link != "" {
								go apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Req.Bearer)
							}
							if utils.Con.SendWebhook {
								go utils.SendWebhook(Req.Info.Name, Req.Bearer)
							}
							fmt.Println(utils.Logo(fmt.Sprintf("[%v] Succesfully sniped! - %v", Req.Info.Name, Req.Email)))
						default:
							Found := 0
							for _, n := range Accs {
								if !n.Taken {
									Found++
								}
							}
							if Found == 0 {
								break Exit
							}
							time.Sleep(10 * time.Second)
						}
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
		var JSON string = fmt.Sprintf(`{"profileName":"%v"}`, name)
		return fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
	} else {
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: bearer " + bearer + "\r\n"
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

var ReqAmt int

func Snipe(Config apiGO.Info, name string, NameRecvChannel *bool, SnipedSingleIGN *chan apiGO.Details, list bool, names *[]utils.Names, ListName *chan string) {
	Next := time.Now()
Exit:
	for {
		if *NameRecvChannel {
			break Exit
		}
		select {
		case Data := <-*ListName:
			if list {
				for i, name := range *names {
					if strings.EqualFold(Data, name.Name) {
						(*names)[i].Taken = true
					}
				}
			}
		default:
			New := Next.Add(time.Duration(utils.Con.TimeBetweenSleeps) * time.Millisecond)
			for _, Acc := range utils.Bearer.Details {
				if strings.EqualFold(Acc.Email, Config.Email) {
					Config = Acc
					break
				}
			}
			time.Sleep(New.Sub(time.Unix(New.Unix()-5, 0)))
			if proxy := utils.Connect(utils.Proxy.CompRand()); proxy.Alive {
				var Payload string
				if list {
					if Data := (*names)[rand.Intn(len(*names))]; !Data.Taken {
						name = Data.Name
						Payload = ReturnPayload(Config.AccountType, Config.Bearer, Data.Name)
					}
				} else {
					Payload = ReturnPayload(Config.AccountType, Config.Bearer, name)
				}
				fmt.Fprint(proxy.Proxy, Payload)
				time.Sleep(time.Until(Next))
				if Payload != "" && !*NameRecvChannel {
					ReqAmt++
					Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(proxy.Proxy, "\r\n"), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType, Info: Config.Info}
					var Details utils.Status
					switch true {
					case strings.Contains(Req.ResponseDetails.Body, "ALREADY_REGISTERED"):
						Details.Data.Status = "ALREADY_REGISTERED"
					case strings.Contains(Req.ResponseDetails.Body, "NOT_ENTITLED"):
						Details.Data.Status = "NOT_ENTITLED"
						// delete acc from bearers and update config on namechange state.
						var New []apiGO.Info
						for _, acc := range utils.Bearer.Details {
							if !strings.EqualFold(acc.Email, Config.Email) {
								New = append(New, acc)
							}
						}
						utils.Bearer.Details = New
						for i, acc := range utils.Con.Bearers {
							if strings.EqualFold(acc.Email, Config.Email) {
								utils.Con.Bearers[i].NameChange = false
								utils.Con.SaveConfig()
								utils.Con.LoadState()
								break
							}
						}
						fmt.Println(utils.Logo(fmt.Sprintf("[401] Account %v has become invalid.. (no longer is a valid Gamepass account)", Config.Email)))
						return
					case strings.Contains(Req.ResponseDetails.Body, "DUPLICATE"):
						Details.Data.Status = "DUPLICATE"
					case strings.Contains(Req.ResponseDetails.Body, "NOT_ALLOWED"):
						Details.Data.Status = "NOT_ALLOWED"
					default:
						switch Req.ResponseDetails.StatusCode {
						case "429":
							Details.Data.Status = "RATE_LIMITED"
						case "401":
							Details.Data.Status = "UNAUTHORIZED"
						case "200":
							Details.Data.Status = "CLAIMED"
						case "":
							Details.Data.Status = "DEAD_PROXY"
						default:
							Details.Data.Status = "UNKNOWN:" + Req.ResponseDetails.StatusCode
						}
					}
					C := fmt.Sprintf(`[%v] <%v> ~ [%v] {"status":"%v","name":"%v","account_type":"%v"}`, ReqAmt, Req.ResponseDetails.SentAt.Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, Details.Data.Status, name, Config.AccountType)
					fmt.Print(utils.Logo(C), "           \r")
					utils.WriteToLogs(name, C+"\n")
					switch Req.ResponseDetails.StatusCode {
					case "200":
						Req.Info.Name = name
						*SnipedSingleIGN <- Req
						*NameRecvChannel = true
						fmt.Println(utils.Logo(fmt.Sprintf("%v claimed %v @ %v\n", Config.Email, name, Req.ResponseDetails.SentAt)))
					}
				} else if list {
					var found bool
					for _, names := range *names {
						if !names.Taken {
							found = true
							break
						}
					}
					if !found {
						break Exit
					}
				}
			}
			Next = New
		}
	}
}
