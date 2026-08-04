package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/chickeniq/warp-wg-go/pkg/warp"
)

func loadAccount(path string) (*warp.Account, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var acc warp.Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return nil, fmt.Errorf("failed to parse account: %w", err)
	}

	return &acc, nil
}

func registerAccount(ctx context.Context, path string) (*warp.Account, error) {
	acc, err := warp.RegisterAccount(ctx)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write warp account file: %w", err)
	}

	return acc, nil
}

func loadOrRegisterAccount(ctx context.Context, path string) (*warp.Account, error) {
	acc, err := loadAccount(path)
	if err == nil {
		return acc, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return registerAccount(ctx, path)
}
