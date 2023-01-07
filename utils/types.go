package utils

import (
	"crypto/x509"
	"main/webhook"

	"github.com/6uf/apiGO"
)

var (
	Roots  *x509.CertPool = x509.NewCertPool()
	Con    Config
	Proxy  apiGO.Proxys
	Bearer apiGO.MCbearers
	RGB    []string
)

type Names struct {
	Name  string
	Taken bool
}

type Proxies struct {
	IP, Port, User, Password string
}

type Status struct {
	Data struct {
		Status string `json:"status"`
	} `json:"details"`
}

type Config struct {
	Gradient           []Values        `json:"gradient"`
	GC_ReqAmt          int             `json:"amt_reqs_per_gc_acc"`
	MFA_ReqAmt         int             `json:"amt_reqs_per_mfa_acc"`
	UseMethod          bool            `json:"use_method_rlbypass"`
	UseProxyDuringAuth bool            `json:"useproxysduringauth"`
	UseCustomSpread    bool            `json:"use_own_spread_value"`
	Spread             int64           `json:"spread_ms"`
	FirstUse           bool            `json:"firstuse"`
	DownloadedPW       bool            `json:"pwinstalled"`
	UseWebhook         bool            `json:"sendpersonalwhonsnipe"`
	WebhookURL         string          `json:"webhook_url"`
	Webhook            webhook.Web     `json:"webhook_json"`
	SkinChange         Skin            `json:"skin_config"`
	Bearers            []apiGO.Bearers `json:"Bearers"`
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
