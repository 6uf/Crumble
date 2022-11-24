package src

import (
	"time"

	"github.com/6uf/apiGO"
)

var (
	Con    Config
	Proxy  apiGO.Proxys
	Bearer apiGO.MCbearers
)

type Proxies struct {
	IP, Port, User, Password string
}

type Config struct {
	SpreadPerAccount   int64           `json:"spreadforaccounts"`
	SpreadPerSend      int64           `json:"spreadforsends"`
	ReqAmtPerAcc       int64           `json:"reqamt"`
	SkinChange         Skin            `json:"skin_config"`
	Bearers            []apiGO.Bearers `json:"Bearers"`
	UseProxyDuringAuth bool            `json:"useproxysduringauth"`
	DiscordID          string          `json:"id"`
	SendWebhook        bool            `json:"sendwebhook"`
}

type Info struct {
	Bearer       string
	RefreshToken string
	AccessToken  string
	Expires      int
	AccountType  string
	Email        string
	Password     string
	Requests     int
	Info         UserINFO `json:"Info"`
	Error        string
}

type UserINFO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Skin struct {
	Link    string `json:"url"`
	Variant string `json:"variant"`
}

type T struct {
	F []TimeFluc
}

type TimeFluc struct {
	T1   time.Time
	Err1 error
}
