package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/chickeniq/cf-router-go/pkg/cf"
	"github.com/chickeniq/cf-router-go/pkg/forwarder"
)

const (
	defaultAccountPath   = "account.json"
	defaultWireguardPath = "wg.json"
	defaultMTU           = 1280
)

type forwardPair struct {
	listen   string
	upstream string
}

func envOrDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	return value
}

func parseMTU(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultMTU, nil
	}

	mtu, err := strconv.Atoi(value)
	if err != nil || mtu <= 0 {
		return 0, fmt.Errorf("invalid MTU %q", value)
	}

	return mtu, nil
}

func parseForwards(value string) ([]forwardPair, error) {
	var forwards []forwardPair
	for pair := range strings.SplitSeq(value, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		listen, upstream, ok := strings.Cut(pair, "=")
		if !ok {
			return nil, fmt.Errorf("invalid FORWARDS pair %q", pair)
		}

		listen = strings.TrimSpace(listen)
		upstream = strings.TrimSpace(upstream)
		if listen == "" || upstream == "" {
			return nil, fmt.Errorf("invalid FORWARDS pair %q", pair)
		}

		forwards = append(forwards, forwardPair{
			listen:   listen,
			upstream: upstream,
		})
	}

	if len(forwards) == 0 {
		return nil, fmt.Errorf("FORWARDS is empty")
	}

	return forwards, nil
}

func main() {
	accountPath := envOrDefault(os.Getenv("ACCOUNT_PATH"), defaultAccountPath)
	wireguardPath := envOrDefault(os.Getenv("WIREGUARD_PATH"), defaultWireguardPath)

	forwards, err := parseForwards(os.Getenv("FORWARDS"))
	if err != nil {
		log.Fatal(err)
	}

	mtu, err := parseMTU(os.Getenv("MTU"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	setupCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cfg, err := cf.GetWireguardConfig(setupCtx, accountPath, wireguardPath)
	if err != nil {
		log.Fatal(err)
	}

	tun, err := cfg.NewTun(mtu)
	if err != nil {
		log.Fatal(err)
	}

	dev, err := cfg.NewDev(tun.Dev)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	dialer := forwarder.Dialer{
		Timeout:     10 * time.Second,
		DialContext: tun.Net.DialContext,
	}

	errs := make(chan error, len(forwards))
	for _, pair := range forwards {
		log.Printf("%s -> %s", pair.listen, pair.upstream)
		go func() {
			if err := forwarder.ForwardConn(ctx, dialer, pair.listen, pair.upstream); err != nil {
				errs <- fmt.Errorf("forward %s to %s: %w", pair.listen, pair.upstream, err)
			}
		}()
	}

	select {
	case <-ctx.Done():
	case err := <-errs:
		log.Print(err)
	}
}
