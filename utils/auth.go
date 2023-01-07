package utils

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/6uf/apiGO"
	"github.com/6uf/h2"
)

type Proxys_Accs struct {
	Proxy string
	Accs  []apiGO.Info
}

var Accs map[string][]Proxys_Accs = make(map[string][]Proxys_Accs)
var Use_gc, Accamt int
var First_gc bool = true

func AuthAccs() {
	grabDetails()
	if len(Con.Bearers) == 0 {
		fmt.Println(Logo("No Bearers have been found, please check your details."))
		return
	} else {
		checkifValid()
		for _, Accs := range Con.Bearers {
			if Accs.NameChange {
				Bearer.Details = append(Bearer.Details, apiGO.Info{
					Bearer:      Accs.Bearer,
					AccountType: Accs.Type,
					Email:       Accs.Email,
					Password:    Accs.Password,
				})
			}
		}
	}
}

func grabDetails() {

	var AccountsVer []string
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		fmt.Println(Logo("Unable to continue, you have no Accounts added."))
		return
	}

	var wgs sync.WaitGroup

	for _, acc_1 := range AccountsVer {
		if acc := strings.Split(acc_1, ":"); !(len(acc) > 1) {
			wgs.Add(1)
			go func(acc_1 string) {
				if len(Proxy.Proxys) > 0 && Con.UseProxyDuringAuth {
					ip, port, user, pass := "", "", "", ""
					switch data := strings.Split(Proxy.CompRand(), ":"); len(data) {
					case 2:
						ip = data[0]
						port = data[1]
					case 4:
						ip = data[0]
						port = data[1]
						user = data[2]
						pass = data[3]
					}
					Bearer.Details = append(Bearer.Details, apiGO.Info{
						Bearer:      acc_1,
						Email:       acc_1[0:12],
						AccountType: apiGO.IsGC(acc_1, &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass}),
						Expires:     int(time.Now().Add(3600000 * time.Second).Unix()),
					})
				} else {
					Bearer.Details = append(Bearer.Details, apiGO.Info{
						Bearer:      acc_1,
						Email:       acc_1[0:12],
						AccountType: apiGO.IsGC(acc_1, nil),
						Expires:     int(time.Now().Add(3600000 * time.Second).Unix()),
					})
				}
				wgs.Done()
			}(acc_1)
		}
	}

	wgs.Wait()

	CheckDupes(AccountsVer)

	if Con.Bearers == nil {
		var wg sync.WaitGroup
		for _, acc_1 := range AccountsVer {
			if acc := strings.Split(acc_1, ":"); len(acc) > 1 {
				if len(Proxy.Proxys) > 0 && Con.UseProxyDuringAuth {
					wg.Add(1)
					go func() {
						ip, port, user, pass := "", "", "", ""
						switch data := strings.Split(Proxy.CompRand(), ":"); len(data) {
						case 2:
							ip = data[0]
							port = data[1]
						case 4:
							ip = data[0]
							port = data[1]
							user = data[2]
							pass = data[3]
						}
						switch info := apiGO.MS_authentication(acc[0], acc[1], &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass}); true {
						case info.Error != "":
							fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", info.Email, info.Error)))
						case info.Bearer != "" && apiGO.CheckChange(info.Bearer, &h2.ProxyAuth{IP: ip, Port: port, User: user, Password: pass}):
							fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashMessage(info.Email, len(info.Email)/4))))
							Con.Bearers = append(Con.Bearers, apiGO.Bearers{
								Bearer:       info.Bearer,
								AuthInterval: 3600,
								AuthedAt:     time.Now().Unix(),
								Type:         info.AccountType,
								Email:        info.Email,
								Password:     info.Password,
								NameChange:   true,
							})
						default:
							fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", acc[0], acc[1])))
						}

						wg.Done()
					}()
				} else {
					switch info := apiGO.MS_authentication(acc[0], acc[1], nil); true {
					case info.Error != "":
						fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", info.Email, info.Error)))
					case info.Bearer != "" && apiGO.CheckChange(info.Bearer, nil):
						fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashMessage(info.Email, len(info.Email)/4))))
						Con.Bearers = append(Con.Bearers, apiGO.Bearers{
							Bearer:       info.Bearer,
							AuthInterval: 3600,
							AuthedAt:     time.Now().Unix(),
							Type:         info.AccountType,
							Email:        info.Email,
							Password:     info.Password,
							NameChange:   true,
						})
					default:
						fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", acc[0], acc[1])))
					}
				}
			}
		}

		wg.Wait()
	} else if len(Con.Bearers) < len(AccountsVer) {
		var auth []string
		check := make(map[string]bool)
		for _, Acc := range Con.Bearers {
			check[Acc.Email+":"+Acc.Password] = true
		}

		for _, Accs := range AccountsVer {
			if !check[Accs] {
				auth = append(auth, Accs)
			}
		}
		var wg sync.WaitGroup
		for _, acc := range auth {
			if acc := strings.Split(acc, ":"); len(acc) > 1 {
				if len(Proxy.Proxys) > 0 && Con.UseProxyDuringAuth {
					wg.Add(1)
					go func() {
						ip, port, user, pass := "", "", "", ""
						switch data := strings.Split(Proxy.CompRand(), ":"); len(data) {
						case 2:
							ip = data[0]
							port = data[1]
						case 4:
							ip = data[0]
							port = data[1]
							user = data[2]
							pass = data[3]
						}
						switch info := apiGO.MS_authentication(acc[0], acc[1], &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass}); true {
						case info.Error != "":
							fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", info.Email, info.Error)))
						case info.Bearer != "" && apiGO.CheckChange(info.Bearer, &h2.ProxyAuth{IP: ip, Port: port, User: user, Password: pass}):
							fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashMessage(info.Email, len(info.Email)/4))))
							Con.Bearers = append(Con.Bearers, apiGO.Bearers{
								Bearer:       info.Bearer,
								AuthInterval: 3600,
								AuthedAt:     time.Now().Unix(),
								Type:         info.AccountType,
								Email:        info.Email,
								Password:     info.Password,
								NameChange:   true,
							})
						default:
							fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", acc[0], acc[1])))
						}
						wg.Done()
					}()
				} else {
					switch info := apiGO.MS_authentication(acc[0], acc[1], nil); true {
					case info.Error != "":
						fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", info.Email, info.Error)))
					case info.Bearer != "" && apiGO.CheckChange(info.Bearer, nil):
						fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashMessage(info.Email, len(info.Email)/4))))
						Con.Bearers = append(Con.Bearers, apiGO.Bearers{
							Bearer:       info.Bearer,
							AuthInterval: 3600,
							AuthedAt:     time.Now().Unix(),
							Type:         info.AccountType,
							Email:        info.Email,
							Password:     info.Password,
							NameChange:   true,
						})
					default:
						fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", acc[0], acc[1])))
					}
				}
			}
		}

		wg.Wait()
	} else if len(AccountsVer) < len(Con.Bearers) {
		var New []apiGO.Bearers
		for _, Accs := range AccountsVer {
			for _, num := range Con.Bearers {
				if Accs == num.Email+":"+num.Password {
					New = append(New, num)
				}
			}
		}
		Con.Bearers = New
	}

	Con.SaveConfig()
	Con.LoadState()
}

func checkifValid() {
	var reAuth []string
	var wgs sync.WaitGroup
	for _, Accs := range Con.Bearers {
		wgs.Add(1)
		go func(Accs apiGO.Bearers) {
			f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
			f.Header.Set("Authorization", "Bearer "+Accs.Bearer)
			j, _ := http.DefaultClient.Do(f)

			if j.StatusCode == 401 {
				fmt.Println(Logo(fmt.Sprintf("Reauthing %v", Accs.Email)))
				reAuth = append(reAuth, Accs.Email+":"+Accs.Password)
			}
			wgs.Done()
		}(Accs)
	}

	wgs.Wait()

	if len(reAuth) != 0 {
		var wg sync.WaitGroup
		for _, acc := range reAuth {
			if acc := strings.Split(acc, ":"); len(acc) > 1 {
				if len(Proxy.Proxys) > 0 && Con.UseProxyDuringAuth {
					wg.Add(1)
					go func() {
						ip, port, user, pass := "", "", "", ""
						switch data := strings.Split(Proxy.CompRand(), ":"); len(data) {
						case 2:
							ip = data[0]
							port = data[1]
						case 4:
							ip = data[0]
							port = data[1]
							user = data[2]
							pass = data[3]
						}
						switch info := apiGO.MS_authentication(acc[0], acc[1], &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass}); true {
						case info.Error != "":
							fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", info.Email, info.Error)))
						case info.Bearer != "" && apiGO.CheckChange(info.Bearer, &h2.ProxyAuth{IP: ip, Port: port, User: user, Password: pass}):
							fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashMessage(info.Email, len(info.Email)/4))))
							Con.Bearers = append(Con.Bearers, apiGO.Bearers{
								Bearer:       info.Bearer,
								AuthInterval: 3600,
								AuthedAt:     time.Now().Unix(),
								Type:         info.AccountType,
								Email:        info.Email,
								Password:     info.Password,
								NameChange:   true,
							})
						default:
							fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", acc[0], acc[1])))
						}
						wg.Done()
					}()
				} else {
					switch info := apiGO.MS_authentication(acc[0], acc[1], nil); true {
					case info.Error != "":
						fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", info.Email, info.Error)))
					case info.Bearer != "" && apiGO.CheckChange(info.Bearer, nil):
						fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashMessage(info.Email, len(info.Email)/4))))
						Con.Bearers = append(Con.Bearers, apiGO.Bearers{
							Bearer:       info.Bearer,
							AuthInterval: 3600,
							AuthedAt:     time.Now().Unix(),
							Type:         info.AccountType,
							Email:        info.Email,
							Password:     info.Password,
							NameChange:   true,
						})
					default:
						fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", acc[0], acc[1])))
					}
				}
			}
		}
		wg.Wait()
	}

	Con.SaveConfig()
	Con.LoadState()
}

// _diamondburned_#4507 thanks to them for the epic example below.

func CheckDupes(strs []string) []string {
	dedup := strs[:0] // re-use the backing array
	track := make(map[string]bool, len(strs))

	for _, str := range strs {
		if track[str] {
			continue
		}
		dedup = append(dedup, str)
		track[str] = true
	}

	return dedup
}

func CheckAccs() {
	for {
		time.Sleep(10 * time.Second)
		var wg sync.WaitGroup
		for _, acc := range Con.Bearers {
			if time.Now().Unix() > acc.AuthedAt+acc.AuthInterval {
				wg.Add(1)
				go func(acc apiGO.Bearers) {
					ip, port, user, pass := "", "", "", ""
					if len(Proxy.Proxys) > 0 && Con.UseProxyDuringAuth {
						switch data := strings.Split(Proxy.CompRand(), ":"); len(data) {
						case 2:
							ip = data[0]
							port = data[1]
						case 4:
							ip = data[0]
							port = data[1]
							user = data[2]
							pass = data[3]
						}
					}
					defer wg.Done()
					switch info := apiGO.MS_authentication(acc.Email, acc.Password, &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass}); true {
					case info.Bearer != "" && apiGO.CheckChange(info.Bearer, &h2.ProxyAuth{IP: ip, Port: port, User: user, Password: pass}) && info.Error == "":
						for point, bf := range Con.Bearers {
							if strings.EqualFold(bf.Email, info.Email) {
								Con.Bearers[point] = apiGO.Bearers{
									Bearer:       info.Bearer,
									NameChange:   true,
									Type:         info.AccountType,
									Password:     info.Password,
									Email:        info.Email,
									AuthedAt:     time.Now().Unix(),
									AuthInterval: 3600,
								}
								break
							}
						}
						for i, Bearers := range Bearer.Details {
							if strings.EqualFold(Bearers.Email, info.Email) {
								Bearer.Details[i] = info
								break
							}
						}
						var Found bool
						for i, accs := range Accs["Giftcard"] {
							for e, b := range accs.Accs {
								if strings.EqualFold(b.Email, info.Email) {
									Accs["Giftcard"][i].Accs[e] = info
									Found = true
									break
								}
							}
						}
						if !Found {
							for i, accs := range Accs["Microsoft"] {
								for e, b := range accs.Accs {
									if strings.EqualFold(b.Email, info.Email) {
										Accs["Microsoft"][i].Accs[e] = info
										Found = true
										break
									}
								}
							}
						}
					}
				}(acc)
			}
		}

		wg.Wait()

		Con.SaveConfig()
		Con.LoadState()
	}
}
