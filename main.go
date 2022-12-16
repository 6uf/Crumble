package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/utils"
	"main/webhook"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/6uf/StrCmd"
	"github.com/6uf/apiGO"
)

var Username = ""

func TempCalc(interval int) time.Duration {
	return time.Duration(interval/len(utils.Bearer.Details)) * time.Millisecond
}

func BuildWebhook(name, searches, headurl string) []byte {
	new := utils.Con.Webhook
	for i := range new.Embeds {
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{name}", name, -1)
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{searches}", searches, -1)
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{id}", utils.Con.DiscordID, -1)
		new.Embeds[i].Author.Name = strings.Replace(new.Embeds[i].Author.Name, "{name}", name, -1)
		new.Embeds[i].Author.Name = strings.Replace(new.Embeds[i].Author.Name, "{searches}", searches, -1)
		new.Embeds[i].Author.IconURL = strings.Replace(new.Embeds[i].Author.IconURL, "{headurl}", headurl, -1)
		new.Embeds[i].Author.IconURL = strings.Replace(new.Embeds[i].Author.IconURL, "{name}", name, -1)
		new.Embeds[i].Author.URL = strings.Replace(new.Embeds[i].Author.URL, "{headurl}", headurl, -1)
		new.Embeds[i].Author.URL = strings.Replace(new.Embeds[i].Author.URL, "{name}", name, -1)
		new.Embeds[i].URL = strings.Replace(new.Embeds[i].URL, "{name}", name, -1)
		new.Embeds[i].Footer.Text = strings.Replace(new.Embeds[i].Footer.Text, "{name}", name, -1)
		new.Embeds[i].Footer.Text = strings.Replace(new.Embeds[i].Footer.Text, "{searches}", searches, -1)
		new.Embeds[i].Footer.IconURL = strings.Replace(new.Embeds[i].Footer.IconURL, "{name}", name, -1)
		new.Embeds[i].Footer.IconURL = strings.Replace(new.Embeds[i].Footer.IconURL, "{headurl}", headurl, -1)
	}
	json, _ := json.Marshal(new)
	return json
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
		fmt.Print(utils.Logo("Use proxys for authentication? : [YES/NO] > "))
		var ProxyAuth string
		fmt.Scan(&ProxyAuth)
		utils.Con.FirstUse = false
		utils.Con.UseProxyDuringAuth = strings.Contains(strings.ToLower(ProxyAuth), "y")
		fmt.Print(utils.Logo("Send to webhook once sniped? : [YES/NO] > "))
		var WebhookSend string
		fmt.Scan(&WebhookSend)
		utils.Con.SendWebhook = strings.Contains(strings.ToLower(WebhookSend), "y")
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
	if file_name := "proxys.txt"; utils.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	utils.Proxy.GetProxys(false, nil)
	utils.Proxy.Setup()
	utils.AuthAccs()
	go utils.CheckAccs()
}

func main() {
	app := StrCmd.App{
		Version:        "v1.3.00-CR",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"webhook": {
				Subcommand: map[string]StrCmd.SubCmd{
					"test": {
						Action: func() {
							_, _, _, searches := utils.GetDroptimes("test")
							err, ok := webhook.Webhook(utils.Con.WebhookURL, BuildWebhook("test", searches, utils.GetHeadUrl("test")))
							if err != nil {
								fmt.Println(utils.Logo(err.Error()))
							} else if ok {
								fmt.Println(utils.Logo("Succesfully sent personal webhook!"))
							}
						},
					},
				},
			},
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					cl, name, Changed, c := false, StrCmd.String("-u"), false, make(chan os.Signal, 1)
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
								if start != 0 && end != 0 && time.Now().After(time.Unix(end, 0)) {
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
Proxies    ~ %v
Account(s) ~ %v
Searches   ~ %v
Status     ~ %v
Start      ~ %v
End        ~ %v

`, name, len(utils.Proxy.Proxys), len(utils.Bearer.Details), searches, status, time.Unix(start, 0), time.Unix(end, 0))))
					var Accs_value []struct {
						Proxy   string
						Account apiGO.Info
					}
					for i, proxy := range utils.Proxy.Proxys {
						if i < len(utils.Bearer.Details) {
							Accs_value = append(Accs_value, struct {
								Proxy   string
								Account apiGO.Info
							}{Proxy: proxy, Account: utils.Bearer.Details[i]})
						}
					}
					for _, Config := range Accs_value {
						Spread, SleepLength := time.Millisecond, time.Millisecond
						if utils.Con.UseCustomSpread {
							Spread = time.Duration(utils.Con.Spread) * time.Millisecond
							if Config.Account.AccountType == "Giftcard" {
								SleepLength = time.Millisecond * 15000
							} else {
								SleepLength = time.Millisecond * 10000
							}
						} else {
							if Config.Account.AccountType == "Giftcard" {
								Spread = TempCalc(15000)
								SleepLength = time.Millisecond * 15000
							} else {
								Spread = TempCalc(10000)
								SleepLength = time.Millisecond * 10000
							}
						}
						go Snipe(Config.Account, SleepLength, name, &Changed, Config.Proxy)
						time.Sleep(Spread)
					}
				Exit:
					for {
						if cl {
							ReqAmt = 0
							fmt.Println()
							fmt.Println(utils.Logo(name + " Has dropped."))
							signal.Stop(c)
							break Exit
						}
						time.Sleep(1 * time.Second)
					}
				},
				Args: []string{
					"-u",
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
var LastReq time.Time

func Snipe(Config apiGO.Info, Spread time.Duration, name string, NameRecvChannel *bool, proxy string) {
	Next := time.Now()
Exit:
	for {
		switch true {
		case *NameRecvChannel:
			break Exit
		default:
			New := Next.Add(Spread)
			LastReq = time.Now()
			for _, Acc := range utils.Bearer.Details {
				if strings.EqualFold(Acc.Email, Config.Email) {
					Config = Acc
					break
				}
			}
			time.Sleep(New.Sub(time.Unix(New.Unix()-5, 0)))
			if proxy := utils.Connect(proxy); proxy.Alive {
				var Payload string = ReturnPayload(Config.AccountType, Config.Bearer, name)
				fmt.Fprint(proxy.Proxy, Payload)
				time.Sleep(time.Until(Next))
				if Payload != "" && !*NameRecvChannel {
					ReqAmt++
					Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(proxy.Proxy, "\r\n"), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType}
					var Details utils.Status
					switch true {
					case strings.Contains(Req.ResponseDetails.Body, "ALREADY_REGISTERED"):
						Details.Data.Status = "ALREADY_REGISTERED"
						UpdateConfig(Config.Email)
						fmt.Println(utils.Logo(fmt.Sprintf("[401] %v cannot name change anymore!", Config.Email)))
						return
					case strings.Contains(Req.ResponseDetails.Body, "NOT_ENTITLED"):
						Details.Data.Status = "NOT_ENTITLED"
						UpdateConfig(Config.Email)
						fmt.Println(utils.Logo(fmt.Sprintf("[401] Account %v has become invalid.. (no longer is a valid Gamepass account)", Config.Email)))
						return
					case strings.Contains(Req.ResponseDetails.Body, "DUPLICATE"):
						Details.Data.Status = "DUPLICATE"
					case strings.Contains(Req.ResponseDetails.Body, "NOT_ALLOWED"):
						Details.Data.Status = "NOT_ALLOWED"
						fmt.Println(utils.Logo(name + " is currently blocked!"))
						return
					default:
						switch Req.ResponseDetails.StatusCode {
						case "429":
							Details.Data.Status = "RATE_LIMITED"
						case "401":
							Details.Data.Status = "UNAUTHORIZED"
						case "200":
							Details.Data.Status = "CLAIMED"
							ReqAmt = 0
							if utils.Con.SkinChange.Link != "" {
								go apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Req.Bearer)
							}
							if utils.Con.SendWebhook {
								go utils.SendWebhook(name, Req.Bearer)
							}
							if utils.Con.UseWebhook {
								go func() {
									_, _, _, searches := utils.GetDroptimes(name)
									err, ok := webhook.Webhook(utils.Con.WebhookURL, BuildWebhook(name, searches, utils.GetHeadUrl(name)))
									if err != nil {
										fmt.Println(utils.Logo(err.Error()))
									} else if ok {
										fmt.Println(utils.Logo("Succesfully sent personal webhook!"))
									}
								}()
							}
							fmt.Println(utils.Logo(fmt.Sprintf("%v claimed %v @ %v\n", Config.Email, name, Req.ResponseDetails.SentAt)))
							*NameRecvChannel = true
						case "":
							Details.Data.Status = "DEAD_PROXY"
						default:
							Details.Data.Status = "UNKNOWN:" + Req.ResponseDetails.StatusCode
						}
					}
					C := fmt.Sprintf(`[%v] %v <%v> ~ [%v] {"status":"%v"} ms since last req %v`, ReqAmt, name, Req.ResponseDetails.SentAt.Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, Details.Data.Status, time.Since(LastReq))
					fmt.Print(utils.Logo(C), "           \r")
				}
			}
			Next = New
		}
	}
}

func UpdateConfig(Email string) {
	var New []apiGO.Info
	for _, acc := range utils.Bearer.Details {
		if !strings.EqualFold(acc.Email, Email) {
			New = append(New, acc)
		}
	}
	utils.Bearer.Details = New
	for i, acc := range utils.Con.Bearers {
		if strings.EqualFold(acc.Email, Email) {
			utils.Con.Bearers[i].NameChange = false
			utils.Con.SaveConfig()
			utils.Con.LoadState()
			break
		}
	}
	accs, _ := os.ReadFile("accounts.txt")
	Scanner := bufio.NewScanner(bytes.NewBuffer(accs))
	var N []string
	for Scanner.Scan() {
		if Text := Scanner.Text(); !strings.EqualFold(Email, strings.Split(Text, ":")[0]) {
			N = append(N, Text)
		}
	}
	os.WriteFile("accounts.txt", []byte(strings.Join(N, "\n")), 0644)
}
