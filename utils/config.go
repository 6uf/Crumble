package utils

import (
	"encoding/json"
	"io"
	"main/webhook"
	"os"
	"strings"
)

var (
	ProxyByte = []byte(`
	-- GlobalSign Root R2, valid until Dec 15, 2021
	-----BEGIN CERTIFICATE-----
	MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
	A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
	Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
	MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
	A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
	hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
	v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
	eoLrvozps6yWJQeXSpkqBy+0Hne/ig+1AnwblrjFuTosvNYSuetZfeLQBoZfXklq
	tTleiDTsvHgMCJiEbKjNS7SgfQx5TfC4LcshytVsW33hoCmEofnTlEnLJGKRILzd
	C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
	zq+r1feqCapgvdzZX99yqWATXgAByUr6P6TqBwMhAo6CygPCm48CAwEAAaOBnDCB
	mTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUm+IH
	V2ccHsBqBt5ZtJot39wZhi4wNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC5n
	bG9iYWxzaWduLm5ldC9yb290LXIyLmNybDAfBgNVHSMEGDAWgBSb4gdXZxwewGoG
	3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
	J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
	291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
	ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
	AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
	TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
	-----END CERTIFICATE-----`)
)

func (s *Config) ToJson() (Data []byte, err error) {
	return json.MarshalIndent(s, "", "  ")
}

func (config *Config) SaveConfig() {
	if Json, err := config.ToJson(); err == nil {
		WriteFile("config.json", string(Json))
	}
}

func (s *Config) LoadState() {
	data, err := ReadFile("config.json")
	if err != nil {
		s.LoadFromFile()
		s.FirstUse = true
		s.SkinChange.Variant = "slim"
		s.SkinChange.Link = "https://textures.minecraft.net/texture/516accb84322ca168a8cd06b4d8cc28e08b31cb0555eee01b64f9175cefe7b75"
		s.Gradient = []Values{{R: "125", G: "110", B: "221"}, {R: "90%", G: "45%", B: "97%"}}
		s.Webhook = webhook.Web{Embeds: []webhook.Embeds{{Description: "<@{id}> has succesfully sniped {namme} with {searches} searches!", URL: "https://namemc.com/profile/{name}", Color: 5814783, Author: webhook.Author{Name: "{name}", URL: "{headurl}", IconURL: "{headurl}"}, Footer: webhook.Footer{Text: "", IconURL: ""}, Fields: []webhook.Fields{{Name: "Example", Value: "Example bio!", Inline: false}}}}}
		s.SaveConfig()
		return
	}

	json.Unmarshal([]byte(data), s)
	s.LoadFromFile()
}

func (c *Config) LoadFromFile() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		jsonFile, _ = os.Create("config.json")
	}

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func WriteFile(path string, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func HashMessage(s string, l int) string { return s[:l] + strings.Repeat("*", len(s[l:])) }
