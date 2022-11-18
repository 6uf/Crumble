package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"main/src"

	"github.com/6uf/StrCmd"
	"github.com/6uf/apiGO"
)

func init() {
	apiGO.Clear()
	fmt.Print(`___________________ 
\_   ___ \______   \
/    \  \/|       _/
\     \___|    |   \
 \______  /____|_  / 
    	\/       \/
`)
	src.Con.LoadState()
	// Checks if a file doesnt exist and creates it.
	if file_name := "accounts.txt"; src.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if file_name := "names.txt"; src.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if file_name := "proxys.txt"; src.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if file_name := "config.json"; src.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.MkdirAll("logs/names", 0755)
	}
	src.Proxy.GetProxys(false, nil)
	src.Proxy.Setup()
	src.AuthAccs()
	go src.CheckAccs()
	if src.Con.DiscordID == "" {
		fmt.Print("Please enter your discord id for identification: ")
		fmt.Scan(&src.Con.DiscordID)
		src.Con.SaveConfig()
		src.Con.LoadState()
	}
	if len(src.Bearer.Details) == 0 {
		fmt.Println("No Bearers Detected..")
		os.Exit(0)
	}
	if src.Con.SpreadPerAccount == 0 && src.Con.SpreadPerSend == 0 {
		fmt.Println("Spread and Send delay cannot be put to 0, this can cause wifi lag or cpu overload!\nIf you think this is a error, please double check your config.json file to make sure the proper values are given.")
		os.Exit(0)
	}
	fmt.Println()
}

func GetDiscordUsername(ID string) string {
	resp, err := http.Get("https://buxflip.com/data/discord/" + ID)
	if err != nil {
		return "Unknown"
	} else {
		if resp.StatusCode == 429 {
			return "Unknown - Rate Limited Via API"
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

func main() {
	app := StrCmd.App{
		Version:        "v3.1.0",
		AppDescription: "Welcome to Crumble!!!",
		Commands: map[string]StrCmd.Command{
			"config": {
				Action: func() {
					if data := StrCmd.Int("-accountspread"); data != 0 {
						src.Con.SpreadPerAccount = int64(data)
					}
					if data := StrCmd.Int("-sendspread"); data != 0 {
						src.Con.SpreadPerSend = int64(data)
					}
					if data := StrCmd.String("-skinlink"); data != "" {
						src.Con.SkinChange.Link = data
					}
					src.Con.SaveConfig()
					src.Con.LoadState()
				},
				Args: []string{
					"-accountspread",
					"-sendspread",
					"-skinlink",
				},
			},
			"name": {
				Description: "The name commands uses the -u arg to take in a name value to pass through into the turbo.",
				Args: []string{
					"-u",
				},
				Action: func() { src.SnipeDefault(StrCmd.String("-u")) },
				Subcommand: map[string]StrCmd.SubCmd{
					"list": {
						Description: "the -l arg takes in a string, format the names like this: name1-name2-name3",
						Args: []string{
							"-l",
						},
						Action: func() {
							names := []string{}
							if line := StrCmd.String("-l"); line != "" {
								names = strings.Split(line, "-")
							} else {
								file, err := os.Open("names.txt")
								if err == nil {
									scanner := bufio.NewScanner(file)
									for scanner.Scan() {
										if name := scanner.Text(); name != "" {
											names = append(names, scanner.Text())
										}
									}
								}
							}
							fmt.Printf("Preparing to snipe %v !\n", names)
							if len(names) != 0 {
								for _, name := range names {
									src.SnipeDefault(name)
								}
							}
						},
					},
				},
			},
		},
	}
	app.Run(fmt.Sprintf("@%v/root: ", GetDiscordUsername(src.Con.DiscordID)))
}
