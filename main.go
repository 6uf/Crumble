package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"main/utils"
	"main/webhook"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/6uf/StrCmd"
	"github.com/6uf/apiGO"
	"github.com/bwmarrin/discordgo"
	"github.com/playwright-community/playwright-go"
)

func TempCalc(interval, reqamt, accamt int) time.Duration {
	return time.Duration((interval*reqamt)/accamt) * time.Millisecond
}

func BuildWebhook(name, searches, headurl string) ([]byte, webhook.Web) {
	new := utils.Con.Webhook
	for i := range new.Embeds {
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{name}", name, -1)
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{searches}", searches, -1)
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
			field.Value = strings.Replace(field.Value, "{headurl}", headurl, -1)
			field.Value = strings.Replace(field.Value, "{searches}", searches, -1)
			field.Value = strings.Replace(field.Value, "{name}", name, -1)
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
	if !utils.Con.DownloadedPW {
		if err := playwright.Install(&playwright.RunOptions{Verbose: true}); err == nil {
			utils.Con.DownloadedPW = true
			utils.Con.SaveConfig()
			utils.Con.LoadState()
		}
	}

	fmt.Print(utils.Logo(`_________                        ______ ______     
__  ____/__________  ________ ______  /____  /____ 
_  /    __  ___/  / / /_  __ '__ \_  __ \_  /_  _ \
/ /___  _  /   / /_/ /_  / / / / /  /_/ /  / /  __/
\____/  /_/    \__,_/ /_/ /_/ /_//_.___//_/  \___/ 

`))
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
	var use_proxy int

	for _, bearer := range utils.Bearer.Details {
		if use_proxy >= len(utils.Proxy.Proxys) && len(utils.Proxy.Proxys) < len(utils.Bearer.Details) {
			break
		}
		switch bearer.AccountType {
		case "Microsoft":
			utils.Accs["Microsoft"] = append(utils.Accs["Microsoft"], utils.Proxys_Accs{Proxy: utils.Proxy.Proxys[use_proxy], Accs: []apiGO.Info{bearer}})
			utils.Accamt++
		case "Giftcard":
			if utils.First_gc {
				utils.Accs["Giftcard"] = []utils.Proxys_Accs{{Proxy: utils.Proxy.Proxys[use_proxy]}}
				utils.First_gc = false
				use_proxy++
			}
			if len(utils.Accs["Giftcard"][utils.Use_gc].Accs) != 5/utils.Con.GC_ReqAmt {
				utils.Accs["Giftcard"][utils.Use_gc].Accs = append(utils.Accs["Giftcard"][utils.Use_gc].Accs, bearer)
				utils.Accamt++
			} else {
				utils.Use_gc++
				utils.Accamt++
				utils.Accs["Giftcard"] = append(utils.Accs["Giftcard"], utils.Proxys_Accs{Proxy: utils.Proxy.Proxys[use_proxy], Accs: []apiGO.Info{bearer}})
				use_proxy++
			}
		}
	}

	fmt.Print(utils.Logo(fmt.Sprintf(`i Accounts Loaded  > <%v>
i Proxies Loaded   > <%v>
i Accounts in use  > <%v>
i Proxys in use    > <%v>
i Accounts Details:
 - GC's Per Proxy  > <%v>
 - MFA's Per Proxy > <%v>
 - Req per GC      > <%v>
 - Req per MFA     > <%v>

`,
		len(utils.Bearer.Details),
		len(utils.Proxy.Proxys),
		utils.Accamt,
		use_proxy,
		5/utils.Con.GC_ReqAmt,
		1,
		utils.Con.GC_ReqAmt,
		utils.Con.MFA_ReqAmt)))
}

func main() {
	app := StrCmd.App{
		Version:        "v1.3.00-CR",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"reload": {
				Action: func() {
					if pw, err := playwright.Run(); err == nil {
						if browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
							Channel:  &[]string{"chrome"}[0],
							Headless: &[]bool{true}[0],
						}); err == nil {
							var wg sync.WaitGroup
							for i, acc := range utils.Bearer.Details {
								wg.Add(1)
								go func(acc apiGO.Info, i int) {
									defer wg.Done()
									if page, err := browser.NewPage(playwright.BrowserNewContextOptions{
										UserAgent: &[]string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"}[0],
									}); err == nil {
										if r, err := page.Goto("https://www.minecraft.net/en-us/login"); err == nil {
											b, _ := page.Content()
											if strings.Contains(b, `You don't have permission to access "http://www.minecraft.net/en-us/login" on this server.`) {
												fmt.Println(utils.Logo(fmt.Sprintf("<%v> %v Currently Ratelimited", i, r.Status())))
											} else {
												if err := page.Click("#main-content > div.page-section.page-section--first.site-content--hide-footer.bg-img-height.bg-globe.d-flex.align-items-center > div > div > div > div.bg-white.py-4 > div:nth-child(1) > div > a"); err != nil {
													for i := 0; i < 5; i++ {
														time.Sleep(30 * time.Second)
														if err := page.Click("#main-content > div.page-section.page-section--first.site-content--hide-footer.bg-img-height.bg-globe.d-flex.align-items-center > div > div > div > div.bg-white.py-4 > div:nth-child(1) > div > a"); err == nil {
															break
														}
													}
												} else {
													time.Sleep(3 * time.Second)
													page.Fill("#i0116", acc.Email)
													fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Added email to input.. "+acc.Email, i)))
													time.Sleep(1 * time.Second)
													page.Click("#idSIButton9")
													fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Submitted email.. "+acc.Email, i)))
													time.Sleep(1 * time.Second)
													page.Fill("#i0118", acc.Password)
													fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Added password to input.. "+acc.Email, i)))
													time.Sleep(1 * time.Second)
													page.Click("#idSIButton9")
													fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Submitted password.. "+acc.Email, i)))

													time.Sleep(2 * time.Second)
													if err := page.Click("#iShowSkip"); err == nil {
														fmt.Println(utils.Logo("i Clicked skip for 7 days.. " + acc.Email))
														time.Sleep(1 * time.Second)
														page.Click("#idBtn_Back")
														time.Sleep(1 * time.Second)
														page.Click("#iCancel")
														time.Sleep(1 * time.Second)
														page.Click("#mc-globalhead__nav-login-dropdown-1 > li:nth-child(2) > a")
														fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Logged out of minecraft.. "+acc.Email, i)))
														time.Sleep(3 * time.Second)
														page.Close()
														fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Succesfully reloaded "+acc.Email, i)))
													} else {
														time.Sleep(1 * time.Second)
														page.Click("#idBtn_Back")
														time.Sleep(1 * time.Second)
														page.Click("#iCancel")
														time.Sleep(1 * time.Second)
														page.Click("#mc-globalhead__nav-login-dropdown-1 > li:nth-child(2) > a")
														fmt.Println(utils.Logo("i Logged out of minecraft.. " + acc.Email))
														time.Sleep(3 * time.Second)
														page.Close()
														fmt.Println(utils.Logo("Succesfully reloaded " + acc.Email))
													}
												}
											}

										} else {
											fmt.Println(err)
										}
									} else {
										fmt.Println(err)
									}
								}(acc, i)
								time.Sleep(3 * time.Second)
							}
							wg.Wait()
						}
					}
				},
			},
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					if len(utils.Con.Bearers) == 0 && len(utils.Bearer.Details) == 0 {
						return
					}
					cl, name, Changed, Use, start, end, Force := false, StrCmd.String("-u"), false, "", 0, 0, StrCmd.Bool("--force")
					if !Force {
						fmt.Println(utils.Logo("Timestamp to Unix: [https://www.epochconverter.com/] (make sure to remove the • on the namemc timestamp!)"))
						fmt.Print(utils.Logo("Use your own unix timestamps [y/n]: "))
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

					if utils.Con.UseMethod {
						for e, p := range utils.Proxy.Proxys {
							if len(utils.Bearer.Details) == e {
								break
							} else {
								Spread, Amt := time.Millisecond, time.Millisecond
								if utils.Bearer.Details[e].AccountType == "Giftcard" {
									Amt = 60060 * time.Millisecond
									if utils.Con.UseCustomSpread {
										Spread = time.Duration(utils.Con.Spread) * time.Millisecond
									} else {
										Spread = TempCalc(15050, 15, len(utils.Bearer.Details))
									}
								} else {
									Amt = 10050 * time.Millisecond
									if utils.Con.UseCustomSpread {
										Spread = time.Duration(utils.Con.Spread) * time.Millisecond
									} else {
										Spread = TempCalc(10050, utils.Con.MFA_ReqAmt, len(utils.Bearer.Details))
									}
								}
								go Snipe(utils.Bearer.Details[e], Amt, name, &Changed, p)
								time.Sleep(Spread)
							}
						}
					} else {
						go func() {
							for _, Acc := range utils.Accs["Giftcard"] {
								Spread := time.Millisecond
								if utils.Con.UseCustomSpread {
									Spread = time.Duration(utils.Con.Spread) * time.Millisecond
								} else {
									Spread = TempCalc(15050, utils.Con.GC_ReqAmt, utils.Accamt)
								}
								for _, data := range Acc.Accs {
									go Snipe(data, time.Duration(15050*utils.Con.GC_ReqAmt)*time.Millisecond, name, &Changed, Acc.Proxy)
									time.Sleep(Spread)
								}
							}
						}()
						go func() {
							for _, Acc := range utils.Accs["Microsoft"] {
								Spread := time.Millisecond
								if utils.Con.UseCustomSpread {
									Spread = time.Duration(utils.Con.Spread) * time.Millisecond
								} else {
									Spread = TempCalc(10050, utils.Con.MFA_ReqAmt, utils.Accamt)
								}
								for _, data := range Acc.Accs {
									go Snipe(data, time.Duration(10050*utils.Con.MFA_ReqAmt)*time.Millisecond, name, &Changed, Acc.Proxy)
									time.Sleep(Spread)
								}
							}
						}()
					}

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

func HashEmailClean(email string) string {
	e := strings.Split(email, "@")[0] // stfu
	return strings.Repeat("⋅", 3) + e[len(e)-4:]
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

func ReturnPayload(acc, bearer, name string) string {
	if acc == "Giftcard" {
		var JSON string = fmt.Sprintf(`{"profileName":"%v"}`, name)
		if utils.Con.UseMethod {
			return fmt.Sprintf("POST https://minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net/minecraft/profile// HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
		} else {
			return fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
		}
	} else {
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: Bearer " + bearer + "\r\n"
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
			if Proxy, ok := utils.Connect(proxy); ok {
				var Payload string = ReturnPayload(Config.AccountType, Config.Bearer, name)
				time.Sleep(time.Until(Next))
				if Payload != "" && !*NameRecvChannel {
					reqamt := 1
					if utils.Con.UseMethod && Config.AccountType == "Giftcard" {
						reqamt = 15
					} else {
						switch Config.AccountType {
						case "Giftcard":
							reqamt = utils.Con.GC_ReqAmt
						case "Microsoft":
							reqamt = utils.Con.MFA_ReqAmt
						}
					}
					var wg sync.WaitGroup
					for i := 0; i < reqamt; i++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(Proxy, Payload), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType}
							var Status string
							switch true {
							case strings.Contains(Req.ResponseDetails.Body, "ALREADY_REGISTERED"):
								Status = "unusable_account"
							case strings.Contains(Req.ResponseDetails.Body, "NOT_ENTITLED"):
								Status = "unusable_account"
								/*
									var ip, port, user, pass string
									switch p := strings.Split(proxy, ":"); len(p) {
									case 2:
										ip, port = p[0], p[1]
									case 4:
										ip, port, user, pass = p[0], p[1], p[2], p[3]
									}
									ReloadAcc(Config.Email, Config.Password, ip, port, user, pass)
								*/
							default:
								switch Req.ResponseDetails.StatusCode {
								case "429":
									//proxy = utils.Proxy.CompRand()
								case "200":
									Status = "claimed"
									if utils.Con.SkinChange.Link != "" {
										go apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Config.Bearer)
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
									fmt.Println(utils.Logo(fmt.Sprintf("⑇ %v claimed %v @ %v\n", Config.Email, name, time.Now().Format("05.0000"))))
									*NameRecvChannel = true
									return
								}
							}
							if Status == "" {
								Status = "unavailable"
							}
							fmt.Println(utils.Logo(fmt.Sprintf(`✗ <%v> [%v] %v ⑇ %v %v ‥ %v`, time.Now().Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, name, Status, HashEmailClean(Config.Email), strings.Split(proxy, ":")[0])))
						}()
					}
					wg.Wait()
				}
			}
			Next = New
		}
	}
}

func ReloadAcc(email, password, ip, port, user, pass string) {
	if pw, err := playwright.Run(); err == nil {
		if browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Channel:  &[]string{"chrome"}[0],
			Headless: &[]bool{true}[0],
			Proxy: &playwright.BrowserTypeLaunchOptionsProxy{
				Server:   playwright.String(fmt.Sprintf("http://%v:%v", ip, port)),
				Username: playwright.String(user),
				Password: playwright.String(pass),
			},
		}); err == nil {
			if page, err := browser.NewPage(playwright.BrowserNewContextOptions{
				UserAgent: &[]string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"}[0],
			}); err == nil {
				if _, err := page.Goto("https://www.minecraft.net/en-us/login"); err == nil {
					b, _ := page.Content()
					if !strings.Contains(b, `You don't have permission to access "http://www.minecraft.net/en-us/login" on this server.`) {
						if err := page.Click("#main-content > div.page-section.page-section--first.site-content--hide-footer.bg-img-height.bg-globe.d-flex.align-items-center > div > div > div > div.bg-white.py-4 > div:nth-child(1) > div > a"); err != nil {
							for i := 0; i < 5; i++ {
								time.Sleep(30 * time.Second)
								if err := page.Click("#main-content > div.page-section.page-section--first.site-content--hide-footer.bg-img-height.bg-globe.d-flex.align-items-center > div > div > div > div.bg-white.py-4 > div:nth-child(1) > div > a"); err == nil {
									break
								}
							}
						} else {
							time.Sleep(3 * time.Second)
							page.Fill("#i0116", email)
							time.Sleep(1 * time.Second)
							page.Click("#idSIButton9")
							time.Sleep(1 * time.Second)
							page.Fill("#i0118", password)
							time.Sleep(1 * time.Second)
							page.Click("#idSIButton9")

							time.Sleep(2 * time.Second)
							if err := page.Click("#iShowSkip"); err == nil {
								time.Sleep(1 * time.Second)
								page.Click("#idBtn_Back")
								time.Sleep(1 * time.Second)
								page.Click("#iCancel")
								time.Sleep(1 * time.Second)
								page.Click("#mc-globalhead__nav-login-dropdown-1 > li:nth-child(2) > a")
								time.Sleep(3 * time.Second)
								page.Close()
							} else {
								time.Sleep(1 * time.Second)
								page.Click("#idBtn_Back")
								time.Sleep(1 * time.Second)
								page.Click("#iCancel")
								time.Sleep(1 * time.Second)
								page.Click("#mc-globalhead__nav-login-dropdown-1 > li:nth-child(2) > a")
								time.Sleep(3 * time.Second)
								page.Close()
							}
						}
					}

				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}
