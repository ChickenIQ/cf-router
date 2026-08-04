package tunnel

import (
	"fmt"

	"github.com/chickeniq/warp-wg-go/pkg/wireguard"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

func NewTunnel(w *wireguard.Config) (*Tunnel, error) {
	addrs, err := parseAddrs(w.Interface.Address)
	if err != nil {
		return nil, err
	}

	dns, err := parseAddrs(w.Interface.DNS)
	if err != nil {
		return nil, err
	}

	tdev, tnet, err := netstack.CreateNetTUN(addrs, dns, w.Interface.MTU)
	if err != nil {
		return nil, fmt.Errorf("failed to create tun: %w", err)
	}

	wgdev, err := w.NewDevice(tdev)
	if err != nil {
		return nil, err
	}

	return &Tunnel{
		WGDev: wgdev,
		Dev:   tdev,
		Net:   tnet,
	}, nil
}

func (t *Tunnel) Close() {
	t.WGDev.Close()
}
