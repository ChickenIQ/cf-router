package wireguard

import (
	"fmt"
	"net"
	"strings"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

func (c *Config) NewTun(mtu int) (*Tun, error) {
	addrs, err := parseAddrs(c.Interface.Address)
	if err != nil {
		return nil, err
	}

	dns, err := parseAddrs(c.Interface.DNS)
	if err != nil {
		return nil, err
	}

	tdev, tnet, err := netstack.CreateNetTUN(addrs, dns, mtu)
	if err != nil {
		return nil, fmt.Errorf("failed to create tun")
	}

	return &Tun{
		Dev: tdev,
		Net: tnet,
	}, nil
}

func (c *Config) NewDev(tunDev tun.Device) (*device.Device, error) {
	dev := device.NewDevice(
		tunDev,
		conn.NewDefaultBind(),
		device.NewLogger(device.LogLevelSilent, ""),
	)
	if dev == nil {
		return nil, fmt.Errorf("failed to create device")
	}

	conf, err := c.IpcConfig()
	if err != nil {
		dev.Close()
		return nil, err
	}

	if err := dev.IpcSet(conf); err != nil {
		dev.Close()
		return nil, fmt.Errorf("failed to configure device: %w", err)
	}

	if err := dev.Up(); err != nil {
		dev.Close()
		return nil, fmt.Errorf("failed to bring device up: %w", err)
	}

	return dev, nil
}

func (c *Config) IpcConfig() (string, error) {
	privateKey, err := keyHex(c.Interface.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	parts := []string{
		strings.Join([]string{
			fmt.Sprintf("private_key=%s", privateKey),
			"replace_peers=true",
		}, "\n"),
	}

	for _, peer := range c.Peers {
		publicKey, err := keyHex(peer.PublicKey)
		if err != nil {
			return "", fmt.Errorf("parse peer public key: %w", err)
		}

		allowedIPs, err := parsePrefixes(peer.AllowedIPs)
		if err != nil {
			return "", err
		}

		lines := []string{
			fmt.Sprintf("public_key=%s", publicKey),
			"persistent_keepalive_interval=25",
			"replace_allowed_ips=true",
		}

		for _, allowedIP := range allowedIPs {
			lines = append(lines, fmt.Sprintf("allowed_ip=%s", allowedIP))
		}

		if peer.Endpoint != "" {
			endpoint, err := net.ResolveUDPAddr("udp", peer.Endpoint)
			if err != nil {
				return "", fmt.Errorf("resolve peer endpoint %q: %w", peer.Endpoint, err)
			}

			lines = append(lines, fmt.Sprintf("endpoint=%s", endpoint))
		}

		parts = append(parts, strings.Join(lines, "\n"))
	}

	return strings.Join(parts, "\n") + "\n", nil
}
