package wireguard

import (
	"encoding/hex"
	"fmt"
	"net/netip"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func keyHex(s string) (string, error) {
	key, err := wgtypes.ParseKey(s)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key[:]), nil
}

func parseAddrs(s string) ([]netip.Addr, error) {
	prefixes, err := parsePrefixes(s)
	if err != nil {
		return nil, err
	}

	addrs := make([]netip.Addr, 0, len(prefixes))
	for _, prefix := range prefixes {
		addrs = append(addrs, prefix.Addr())
	}

	return addrs, nil
}

func parsePrefixes(s string) ([]netip.Prefix, error) {
	parts := strings.Split(s, ",")
	prefixes := make([]netip.Prefix, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if !strings.Contains(part, "/") {
			addr, err := netip.ParseAddr(part)
			if err != nil {
				return nil, fmt.Errorf("parse address %q: %w", part, err)
			}

			prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
			continue
		}

		prefix, err := netip.ParsePrefix(part)
		if err != nil {
			return nil, fmt.Errorf("parse prefix %q: %w", part, err)
		}

		prefixes = append(prefixes, prefix)
	}

	return prefixes, nil
}
