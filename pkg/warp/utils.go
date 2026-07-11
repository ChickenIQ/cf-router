package warp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func newBody(pubKey wgtypes.Key) ([]byte, error) {
	return json.Marshal(map[string]string{
		"fcm_token":  "",
		"install_id": "",
		"key":        pubKey.String(),
		"locale":     "en_US",
		"model":      "PC",
		"tos":        time.Now().Format(time.RFC3339Nano),
		"type":       "Android",
	})
}

func sendRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS12, MaxVersion: tls.VersionTLS12},
		ForceAttemptHTTP2: false,
	}}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "okhttp/3.12.1")
	req.Header.Set("CF-Client-Version", "a-6.3-1922")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: %s", res.Status, data)
	}

	return data, nil
}
