package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chickeniq/cf-router/pkg/cloudflare"
	"github.com/chickeniq/cf-router/pkg/forwarder"
)

func envOrDefault(name string, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}

	return val
}

func mustEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("%s must be set", name)
	}

	return val
}

func main() {
	accPath := envOrDefault("ACCOUNT_PATH", "account.json")
	cfgPath := envOrDefault("WIREGUARD_PATH", "wg.json")
	forwards, err := forwarder.ParseAddrs(mustEnv("FORWARDS"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fetchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tun, err := cloudflare.NewTunnel(fetchCtx, accPath, cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer tun.Close()

	dialer := forwarder.Dialer{Context: tun.Net.DialContext, Timeout: 15 * time.Second}
	errs := make(chan error, len(forwards))

	for _, pair := range forwards {
		log.Printf("%s -> %s", pair.Local, pair.Remote)
		go func() {
			if err := forwarder.ForwardConn(ctx, dialer, pair); err != nil {
				errs <- fmt.Errorf("forward %s to %s: %w", pair.Local, pair.Remote, err)
			}
		}()
	}

	select {
	case <-ctx.Done():
	case err := <-errs:
		log.Print(err)
	}
}
