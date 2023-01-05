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
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/6uf/StrCmd"
	"github.com/6uf/apiGO"
	"github.com/bwmarrin/discordgo"
)

func TempCalc(interval int) time.Duration {
	if utils.Con.UseMethod {
		return time.Duration(interval/(len(utils.Bearer.Details)*15)) * time.Millisecond
	}
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
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}
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
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					if len(utils.Con.Bearers) == 0 && len(utils.Bearer.Details) == 0 {
						return
					}
					cl, name, Changed, Use, start, end, Force := false, StrCmd.String("-u"), false, "", 0, 0, StrCmd.Bool("--force")
					if !Force {
						fmt.Println(utils.Logo("Timestamp to Unix: [https://www.epochconverter.com/] (make sure to remove the • on the namemc timestamp!)"))
						fmt.Print(utils.Logo("Use your own unix timestamps: "))
						fmt.Scan(&Use)
						var start, end int64
						if strings.Contains(strings.ToLower(Use), "y") {
							fmt.Print(utils.Logo("Start: "))
							fmt.Scan(&start)
							fmt.Print(utils.Logo("End: "))
							fmt.Scan(&end)
						}
					}

					drop := time.Unix(int64(start), 0)
					for time.Now().Before(drop) {
						fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v                 \r", name, time.Until(drop).Round(time.Second)))))
						time.Sleep(time.Second * 1)
					}
					go func() {
					Exit:
						for {
							if !Force {
								if utils.IsAvailable(name) {
									Changed = true
									cl = true
									break Exit
								}
								if start != 0 && end != 0 && time.Now().After(time.Unix(int64(end), 0)) {
									Changed = true
									cl = true
									break Exit
								}
							}
							time.Sleep(10 * time.Second)
						}
					}()
					fmt.Print(utils.Logo(fmt.Sprintf(`
Name       ~ %v
Proxies    ~ %v
Account(s) ~ %v
Start      ~ %v
End        ~ %v

`, name, len(utils.Proxy.Proxys), len(utils.Bearer.Details), time.Unix(int64(start), 0), time.Unix(int64(end), 0))))
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
								if gc <= 5 {
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
						for _, Acc := range append(Accs["Giftcard"], Accs["Microsoft"]...) {
							for _, data := range Acc.Accs {
								go func(data apiGO.Info, l int) {
									sp := strings.Split(Acc.Proxy, ":")
									var Proxy func(*http.Request) (*url.URL, error)
									if len(sp) > 2 {
										ip, port, user, pass := sp[0], sp[1], sp[2], sp[3]
										Proxy = http.ProxyURL(&url.URL{Scheme: "http", User: url.UserPassword(user, pass), Host: ip + ":" + port})
									} else {
										ip, port := sp[0], sp[1]
										Proxy = http.ProxyURL(&url.URL{Scheme: "http", Host: ip + ":" + port})
									}
									Client := http.Client{
										Transport: &http.Transport{
											Proxy: Proxy,
										},
									}
									Next := time.Now()
								Exit:
									for {
										switch true {
										case Changed:
											break Exit
										default:
											var New time.Time
											if data.AccountType == "Giftcard" {
												if utils.Con.UseMethod {
													New = Next.Add(5 * time.Second)
												} else {
													New = Next.Add(15 * time.Second)
												}
											} else {
												if utils.Con.UseMethod {
													New = Next.Add(3 * time.Second)
												} else {
													New = Next.Add(10 * time.Second)
												}
											}
											for _, Acc := range utils.Bearer.Details {
												if strings.EqualFold(Acc.Email, data.Email) {
													data = Acc
													break
												}
											}
											time.Sleep(New.Sub(time.Unix(New.Unix()-5, 0)))
											var req *http.Request
											switch data.AccountType {
											case "Microsoft":
												req, _ = http.NewRequest("PUT", "https://api.minecraftservices.com/minecraft/profile/name/"+name, nil)
											case "Giftcard":
												var JSON string = fmt.Sprintf(`{"profileName":"%v"}`, name)
												req, _ = http.NewRequest("POST", "https://api.minecraftservices.com/minecraft/profile//", bytes.NewBuffer([]byte(JSON)))
											}
											req.Header.Add("Authorization", "Bearer "+data.Bearer)
											time.Sleep(time.Until(Next))
											reqamt := 1
											if utils.Con.UseMethod && data.AccountType != "Microsoft" {
												reqamt = l
											}
											for i := 0; i < reqamt; i++ {
												go func() {
													if resp, err := Client.Do(req); err == nil {
														body, _ := io.ReadAll(resp.Body)
														var Details utils.Status

														switch true {
														case strings.Contains(string(body), "ALREADY_REGISTERED"):
															Details.Data.Status = "ALREADY_REGISTERED"
														case strings.Contains(string(body), "NOT_ENTITLED"):
															Details.Data.Status = "NOT_ENTITLED"
														case strings.Contains(string(body), "DUPLICATE"):
															Details.Data.Status = "DUPLICATE"
														case strings.Contains(string(body), "NOT_ALLOWED"):
															Details.Data.Status = "NOT_ALLOWED"
															fmt.Println(utils.Logo(name + " is currently blocked!"))
															return
														default:
															switch resp.StatusCode {
															case 429:
																Details.Data.Status = "RATE_LIMITED"
																/*
																	sp := strings.Split(utils.Proxy.CompRand(), ":")
																	var Proxy func(*http.Request) (*url.URL, error)
																	if len(sp) > 2 {
																		ip, port, user, pass := sp[0], sp[1], sp[2], sp[3]
																		Proxy = http.ProxyURL(&url.URL{Scheme: "http", User: url.UserPassword(user, pass), Host: ip + ":" + port})
																	} else {
																		ip, port := sp[0], sp[1]
																		Proxy = http.ProxyURL(&url.URL{Scheme: "http", Host: ip + ":" + port})
																	}
																	Client.Transport = &http.Transport{Proxy: Proxy}
																*/
															case 401:
																Details.Data.Status = "UNAUTHORIZED"
															case 404:
																for i, Acc := range utils.Con.Bearers {
																	if strings.EqualFold(Acc.Email, data.Email) {
																		utils.Con.Bearers[i].Type = "Giftcard"
																		data.AccountType = "Giftcard"
																		utils.Con.SaveConfig()
																		utils.Con.LoadState()
																		break
																	}
																}
															case 200:
																Details.Data.Status = "CLAIMED"
																if utils.Con.SkinChange.Link != "" {
																	go apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), data.Bearer)
																}
																if utils.Con.UseWebhook {
																	go func() {
																		json, _ := BuildWebhook(name, "0", "")
																		err, ok := webhook.Webhook(utils.Con.WebhookURL, json)
																		if err != nil {
																			fmt.Println(utils.Logo(err.Error()))
																		} else if ok {
																			fmt.Println(utils.Logo("Succesfully sent personal webhook!"))
																		}
																	}()
																}
																fmt.Println(utils.Logo(fmt.Sprintf("%v claimed %v @ %v\n", data.Email, name, time.Now().Format("05.0000"))))
																Changed = true
															default:
																if strings.Contains(string(body), "error") && strings.Contains(string(body), "FORBIDDEN") {
																	Details.Data.Status = "CANNOT_NAME_CHANGE:" + resp.Status + ":" + data.Email
																} else {
																	Details.Data.Status = "UNKNOWN:" + resp.Status
																}
															}
														}
														fmt.Println(utils.Logo(fmt.Sprintf(`• %v <%v> ~ [%v] %v <%v:%v>`, name, time.Now().Format("15:04:05.0000"), resp.StatusCode, Details.Data.Status, data.AccountType, data.Email)))
													}
												}()
												time.Sleep(5 * time.Millisecond)
											}

											if data.AccountType == "Giftcard" {
												if utils.Con.UseMethod {
													New = Next.Add(5 * time.Second)
												} else {
													New = Next.Add(15 * time.Second)
												}
											} else {
												if utils.Con.UseMethod {
													New = Next.Add(3 * time.Second)
												} else {
													New = Next.Add(10 * time.Second)
												}
											}

											Next = New
										}
									}
								}(data, 15/len(Acc.Accs))
								if utils.Con.UseCustomSpread {
									time.Sleep(time.Duration(utils.Con.Spread) * time.Millisecond)
								} else {
									if data.AccountType == "Giftcard" {
										time.Sleep(TempCalc(15000))
									} else {
										time.Sleep(TempCalc(10000))
									}
								}
							}
						}
					}()
				Exit:
					for {
						if cl {
							fmt.Println()
							fmt.Println(utils.Logo(name + " has dropped or ctrl-c was pressed."))
							break Exit
						}
						time.Sleep(1 * time.Second)
					}
				},
				Args: []string{
					"-u",
					"--force",
				},
			},
		},
	}
	app.Run(utils.Logo("@Crumble/root: "))
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
