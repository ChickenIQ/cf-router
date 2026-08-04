package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chickeniq/warp-wg-go/pkg/warp"
	"github.com/chickeniq/warp-wg-go/pkg/wireguard"
)

func loadConfig(path string) (*wireguard.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg wireguard.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse wireguard configv: %w", err)
	}

	return &cfg, nil
}

func registerConfig(ctx context.Context, account *warp.Account, path string) (*wireguard.Config, error) {
	cfg, err := account.RegisterWireguard(ctx)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("failed to write wireguard config file: %w", err)
	}

	return cfg, nil
}
