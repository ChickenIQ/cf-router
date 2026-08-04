package forwarder

import (
	"context"
	"fmt"
	"net"
	"strings"
)

func ForwardConn(ctx context.Context, dialer Dialer, addr AddrPair) error {
	l, err := net.Listen("tcp", addr.Local)
	if err != nil {
		return fmt.Errorf("error listening on %s %s", addr.Local, err.Error())
	}
	defer l.Close()

	go func() {
		<-ctx.Done()
		_ = l.Close()
	}()

	for {
		local, err := l.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("error accepting connection %s", err.Error())
		}

		go connect(ctx, dialer, local, addr.Remote)
	}
}

func ParseAddrs(value string) ([]AddrPair, error) {
	var addrs []AddrPair

	for pair := range strings.SplitSeq(value, ",") {
		if pair == "" {
			return nil, fmt.Errorf("empty addr pair")
		}

		listen, upstream, ok := strings.Cut(pair, "=")
		if !ok || listen == "" || upstream == "" {
			return nil, fmt.Errorf("invalid addr pair %q", pair)
		}

		addrs = append(addrs, AddrPair{
			Local:  listen,
			Remote: upstream,
		})
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("empty addr pairs")
	}

	return addrs, nil
}
