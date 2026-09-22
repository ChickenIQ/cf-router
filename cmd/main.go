package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/chickeniq/cf-router/pkg/cloudflare"
	"github.com/chickeniq/cf-router/pkg/env"
	"github.com/chickeniq/cf-router/pkg/forwarder"
)

const timeout = 15 * time.Second

func main() {
	accountPath := env.Default("ACCOUNT_PATH", "account.json")
	configPath := env.Default("CONFIG_PATH", "config.json")
	addrs, err := forwarder.ParseAddrs(env.Must("FORWARDS"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tun, err := cloudflare.NewTunnel(fetchCtx, accountPath, configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer tun.Close()

	dialer := forwarder.Dialer{Context: tun.Net.DialContext, Timeout: timeout}
	for _, pair := range addrs {
		log.Printf("%s -> %s", pair.Local, pair.Remote)
		go func() {
			if err := forwarder.ForwardConn(ctx, dialer, pair); err != nil {
				log.Printf("forward %s to %s: %v", pair.Local, pair.Remote, err)
				stop()
			}
		}()
	}

	<-ctx.Done()
}
