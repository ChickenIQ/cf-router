package forwarder

import (
	"context"
	"net"
	"time"
)

type DialContext func(context.Context, string, string) (net.Conn, error)

type Dialer struct {
	Timeout     time.Duration
	DialContext DialContext
}
