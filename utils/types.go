package utils

import (
	"crypto/x509"
	"time"

	"github.com/6uf/apiGO"
)

var (
	Roots  *x509.CertPool = x509.NewCertPool()
	Con    Config
	Proxy  apiGO.Proxys
	Bearer apiGO.MCbearers
)

type Proxies struct {
	IP, Port, User, Password string
}

type Config struct {
	Gradient           Gradient        `json:"gradient"`
	SkinChange         Skin            `json:"skin_config"`
	UseProxyDuringAuth bool            `json:"useproxysduringauth"`
	DiscordID          string          `json:"id"`
	SendWebhook        bool            `json:"sendwebhook"`
	Bearers            []apiGO.Bearers `json:"Bearers"`
	FirstUse           bool            `json:"firstuse"`
}

type Gradient struct {
	RGB1 Values `json:"rgb"`
	RGB2 Values `json:"rgb2"`
	HSL  Values `json:"hsl"`
}

type Values struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
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
