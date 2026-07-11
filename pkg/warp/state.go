package warp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func LoadAccount(path string) (*Account, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var acc Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return nil, err
	}

	if acc.ID == "" || acc.Token == "" || acc.PrivateKey == "" {
		return nil, fmt.Errorf("invalid warp data")
	}

	return &acc, nil
}

func SaveAccount(path string, acc *Account) error {
	data, err := json.Marshal(acc)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadOrRegister(ctx context.Context, path string) (*Account, error) {
	acc, err := LoadAccount(path)
	if err == nil {
		return acc, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	acc, err = NewAccount(ctx)
	if err != nil {
		return nil, err
	}

	if err := SaveAccount(path, acc); err != nil {
		return nil, err
	}

	return acc, nil
}
