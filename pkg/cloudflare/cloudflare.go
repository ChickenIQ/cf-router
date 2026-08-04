package cloudflare

import (
	"context"
	"errors"
	"os"

	"github.com/chickeniq/cf-router/pkg/tunnel"
	"github.com/chickeniq/warp-wg-go/pkg/wireguard"
)

func LoadOrRegister(ctx context.Context, accountPath, cfgPath string) (*wireguard.Config, error) {
	cfg, err := loadConfig(cfgPath)
	if err == nil {
		return cfg, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	account, err := loadOrRegisterAccount(ctx, accountPath)
	if err != nil {
		return nil, err
	}

	return registerConfig(ctx, account, cfgPath)
}

func NewTunnel(ctx context.Context, accountPath, cfgPath string) (*tunnel.Tunnel, error) {
	cfg, err := LoadOrRegister(ctx, accountPath, cfgPath)
	if err != nil {
		return nil, err
	}

	return tunnel.NewTunnel(cfg)
}
