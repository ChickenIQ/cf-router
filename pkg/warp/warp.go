package warp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chickeniq/cf-router-go/pkg/wireguard"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func NewAccount(ctx context.Context) (*Account, error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	body, err := newBody(key.PublicKey())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfURL+"/reg", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	data, err := sendRequest(req)
	if err != nil {
		return nil, err
	}

	acc := &Account{PrivateKey: key.String()}
	if err := json.Unmarshal(data, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (acc *Account) WireguardConfig(ctx context.Context) (*wireguard.Config, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfURL+"/reg/"+acc.ID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+acc.Token)

	data, err := sendRequest(req)
	if err != nil {
		return nil, err
	}

	var cfg DeviceConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	addr := cfg.Config.Interface.Addresses
	peer := cfg.Config.Peers[0]

	return &wireguard.Config{
		Interface: wireguard.Interface{
			Address:    fmt.Sprintf("%s/32, %s/128", addr.V4, addr.V6),
			PrivateKey: acc.PrivateKey,
			DNS:        wgDNS,
			MTU:        wgMTU,
		},
		Peers: []wireguard.Peer{
			{
				AllowedIPs: wgIPs,
				PublicKey:  peer.PublicKey,
				Endpoint:   peer.Endpoint.Host,
			},
		},
	}, nil
}
