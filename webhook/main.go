package webhook

import (
	"bytes"
	"errors"
	"net/http"
)

func Webhook(url string, webhook []byte) (error, bool) {
	if resp, err := http.Post(url, "application/json", bytes.NewBuffer(webhook)); err != nil {
		return err, false
	} else {
		if resp.StatusCode == 204 {
			return nil, true
		} else {
			return errors.New("[Error] unable to send webhook " + resp.Status), false
		}
	}
}
