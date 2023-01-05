package utils

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/iskaa02/qalam/gradient"
)

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
}

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
}

func IsAvailable(name string) bool {
	resp, err := http.Get("https://account.mojang.com/available/minecraft/" + name)
	if err == nil {
		return resp.StatusCode == 200
	} else {
		return false
	}
}

func Logo(Data string) string {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors(RGB...).
		Mode(gradient.BlendRgb).
		Build()
	return g.Mutline(Data)
}
