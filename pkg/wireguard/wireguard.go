package wireguard

import (
	"fmt"
	"strings"
)

func (c *Config) String() string {
	parts := []string{
		strings.Join([]string{
			"[Interface]",
			fmt.Sprintf("PrivateKey = %s", c.Interface.PrivateKey),
			fmt.Sprintf("Address = %s", c.Interface.Address),
			fmt.Sprintf("DNS = %s", c.Interface.DNS),
			fmt.Sprintf("MTU = %d", c.Interface.MTU),
		}, "\n"),
	}

	for _, peer := range c.Peers {
		parts = append(parts, strings.Join([]string{
			"[Peer]",
			fmt.Sprintf("PublicKey = %s", peer.PublicKey),
			fmt.Sprintf("AllowedIPs = %s", peer.AllowedIPs),
			fmt.Sprintf("Endpoint = %s", peer.Endpoint),
		}, "\n"))
	}

	return strings.Join(parts, "\n\n") + "\n"
}
