package forwarder

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
)

func connect(ctx context.Context, dialer Dialer, local net.Conn, upstreamAddr string) {
	defer local.Close()

	dialCtx, cancel := context.WithTimeout(ctx, dialer.Timeout)
	defer cancel()

	upstream, err := dialer.Context(dialCtx, "tcp", upstreamAddr)
	if err != nil {
		if ctx.Err() == nil {
			log.Printf("error dialing %s %s", upstreamAddr, err.Error())
		}
		return
	}
	defer upstream.Close()

	if err := copyConn(ctx, local, upstream); err != nil && err.Error() != "done" {
		log.Printf("error forwarding connection %s", err.Error())
	}
}

func copyConn(ctx context.Context, from, to net.Conn) error {
	ctx, cancel := context.WithCancel(ctx)
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		io.Copy(from, to)
		cancel()

		return fmt.Errorf("done")
	})

	eg.Go(func() error {
		io.Copy(to, from)
		cancel()

		return fmt.Errorf("done")
	})

	eg.Go(func() error {
		<-ctx.Done()
		from.Close()
		to.Close()

		return fmt.Errorf("done")
	})

	return eg.Wait()
}
