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
	"reflect"
	"strings"
	"time"

	"github.com/6uf/StrCmd"
	"github.com/6uf/apiGO"
	"github.com/bwmarrin/discordgo"
)

var Username = ""

func TempCalc(interval int) time.Duration {
	return time.Duration(interval/len(utils.Bearer.Details)) * time.Millisecond
}

func BuildWebhook(name, searches, headurl string) ([]byte, webhook.Web) {
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
		for e, field := range new.Embeds[i].Fields {
			field.Name = strings.Replace(field.Name, "{headurl}", headurl, -1)
			field.Name = strings.Replace(field.Name, "{searches}", searches, -1)
			field.Name = strings.Replace(field.Name, "{name}", name, -1)
			field.Name = strings.Replace(field.Name, "{id}", utils.Con.DiscordID, -1)
			field.Value = strings.Replace(field.Value, "{headurl}", headurl, -1)
			field.Value = strings.Replace(field.Value, "{searches}", searches, -1)
			field.Value = strings.Replace(field.Value, "{name}", name, -1)
			field.Value = strings.Replace(field.Value, "{id}", utils.Con.DiscordID, -1)
			new.Embeds[i].Fields[e] = field
		}
	}
	json, _ := json.Marshal(new)
	return json, new
}

func ReturnEmbed(name, searches, headurl string) (Data discordgo.MessageSend) {
	_, new := BuildWebhook(name, searches, headurl)
	for _, com := range new.Embeds {
		var Footer discordgo.MessageEmbedFooter
		var Image discordgo.MessageEmbedImage
		var Thumbnail discordgo.MessageEmbedThumbnail
		var Author discordgo.MessageEmbedAuthor

		if !reflect.DeepEqual(com.Footer, webhook.Footer{}) {
			Footer = discordgo.MessageEmbedFooter{
				Text:    com.Footer.Text,
				IconURL: com.Footer.IconURL,
			}
		}
		if !reflect.DeepEqual(com.Image, webhook.Image{}) {
			Image = discordgo.MessageEmbedImage{
				URL: com.Image.URL,
			}
		}
		if !reflect.DeepEqual(com.Thumbnail, webhook.Thumbnail{}) {
			Thumbnail = discordgo.MessageEmbedThumbnail{
				URL: com.Thumbnail.URL,
			}
		}
		if !reflect.DeepEqual(com.Author, webhook.Author{}) {
			Author = discordgo.MessageEmbedAuthor{
				URL:     com.Author.URL,
				Name:    com.Author.Name,
				IconURL: com.Author.IconURL,
			}
		}

		Data.Embeds = append(Data.Embeds, &discordgo.MessageEmbed{
			URL:         com.URL,
			Description: com.Description,
			Color:       com.Color,
			Footer:      &Footer,
			Image:       &Image,
			Thumbnail:   &Thumbnail,
			Author:      &Author,
			Fields:      returnjustfields(com),
		})
	}
	return
}

func returnjustfields(com webhook.Embeds) (Data []*discordgo.MessageEmbedField) {
	for _, c := range com.Fields {
		Data = append(Data, &discordgo.MessageEmbedField{
			Name:   c.Name,
			Value:  c.Value,
			Inline: c.Inline,
		})
	}
	return
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
							json, _ := BuildWebhook("test", searches, utils.GetHeadUrl("test"))
							err, ok := webhook.Webhook(utils.Con.WebhookURL, json)
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
					if len(utils.Con.Bearers) == 0 && len(utils.Bearer.Details) == 0 {
						return
					}
					cl, name, Changed, c := false, StrCmd.String("-u"), false, make(chan os.Signal, 1)
					signal.Notify(c, os.Interrupt)
					start, end, status, searches := utils.GetDroptimes(name)
					drop := time.Unix(start, 0)
					for time.Now().Before(drop) {
						select {
						case <-c:
							signal.Stop(c)
							Changed = true
							cl = true
							return
						default:
							fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v                 \r", name, time.Until(drop).Round(time.Second)))))
							time.Sleep(time.Second * 1)
						}
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
					type Proxys_Accs struct {
						Proxy string
						Accs  []apiGO.Info
					}
					var Accs map[string][]Proxys_Accs = make(map[string][]Proxys_Accs)
					Accs["Giftcard"] = []Proxys_Accs{{Proxy: utils.Proxy.CompRand()}}
					Accs["Microsoft"] = []Proxys_Accs{{Proxy: utils.Proxy.CompRand()}}
					var gc, use_gc int
					for i, proxy := range utils.Proxy.Proxys {
						if i < len(utils.Bearer.Details) {
							switch utils.Bearer.Details[i].AccountType {
							case "Microsoft":
								Accs["Microsoft"] = append(Accs["Microsoft"], Proxys_Accs{Proxy: proxy, Accs: []apiGO.Info{utils.Bearer.Details[i]}})
							case "Giftcard":
								if gc <= 3 {
									gc++
									Accs["Giftcard"][use_gc].Accs = append(Accs["Giftcard"][use_gc].Accs, utils.Bearer.Details[i])
								} else {
									use_gc++
									gc = 0
									Accs["Giftcard"] = append(Accs["Giftcard"], Proxys_Accs{Proxy: proxy, Accs: []apiGO.Info{utils.Bearer.Details[i]}})
								}
							}
						}
					}

					go func() {
						for _, Acc := range Accs["Giftcard"] {
							Spread := time.Millisecond
							if utils.Con.UseCustomSpread {
								Spread = time.Duration(utils.Con.Spread) * time.Millisecond
							} else {
								Spread = TempCalc(15050)
							}
							for _, data := range Acc.Accs {
								go Snipe(data, 15050*time.Millisecond, name, &Changed, Acc.Proxy)
								time.Sleep(Spread)
							}
						}
					}()
					go func() {
						for _, Acc := range Accs["Microsoft"] {
							Spread := time.Millisecond
							if utils.Con.UseCustomSpread {
								Spread = time.Duration(utils.Con.Spread) * time.Millisecond
							} else {
								Spread = TempCalc(10050)
							}
							for _, data := range Acc.Accs {
								go Snipe(data, 10050*time.Millisecond, name, &Changed, Acc.Proxy)
								time.Sleep(Spread)
							}
						}
					}()
				Exit:
					for {
						if cl {
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
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: Bearer " + bearer + "\r\n"
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

func Snipe(Config apiGO.Info, Spread time.Duration, name string, NameRecvChannel *bool, proxy string) {
	Next := time.Now()
Exit:
	for {
		switch true {
		case *NameRecvChannel:
			break Exit
		default:
			New := Next.Add(Spread)
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
							if utils.Con.SkinChange.Link != "" {
								go apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Req.Bearer)
							}
							if utils.Con.SendWebhook {
								go utils.SendWebhook(name, Req.Bearer)
							}
							if utils.Con.UseWebhook {
								go func() {
									_, _, _, searches := utils.GetDroptimes(name)
									json, _ := BuildWebhook(name, searches, utils.GetHeadUrl(name))
									err, ok := webhook.Webhook(utils.Con.WebhookURL, json)
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
							if strings.Contains(Req.ResponseDetails.Body, "error") && strings.Contains(Req.ResponseDetails.Body, "FORBIDDEN") {
								Details.Data.Status = "CANNOT_NAME_CHANGE:" + Req.ResponseDetails.StatusCode + ":" + Req.Email
							} else {
								Details.Data.Status = "UNKNOWN:" + Req.ResponseDetails.StatusCode
							}
						}
					}
					fmt.Println(utils.Logo(fmt.Sprintf(`%v <%v> ~ [%v] %v <%v>`, name, Req.ResponseDetails.SentAt.Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, Details.Data.Status, Req.Type)))
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
