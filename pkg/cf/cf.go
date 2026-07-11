package cf

import (
	"context"
	"errors"
	"os"

	"github.com/chickeniq/cf-router-go/pkg/warp"
	"github.com/chickeniq/cf-router-go/pkg/wireguard"
)

func GetWireguardConfig(ctx context.Context, cfStatePath, wgStatePath string) (*wireguard.Config, error) {
	conf, err := wireguard.LoadConfig(wgStatePath)
	if errors.Is(err, os.ErrNotExist) {
		account, err := warp.LoadOrRegister(ctx, cfStatePath)
		if err != nil {
			return nil, err
		}

		conf, err = account.WireguardConfig(ctx)
		if err != nil {
			return nil, err
		}

		if err := wireguard.SaveConfig(conf, wgStatePath); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return conf, nil
}
