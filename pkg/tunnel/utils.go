package tunnel

import (
	"net/netip"

	"github.com/chickeniq/warp-wg-go/pkg/wireguard"
)

func parseAddrs(s string) ([]netip.Addr, error) {
	prefixes, err := wireguard.ParsePrefixes(s)
	if err != nil {
		return nil, err
	}

	addrs := make([]netip.Addr, 0, len(prefixes))
	for _, prefix := range prefixes {
		addrs = append(addrs, prefix.Addr())
	}

	return addrs, nil
}
