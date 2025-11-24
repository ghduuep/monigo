package monitor

import (
	"context"
	"net"
	"time"
)

func newCustomResolver(dnsServerAddr string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,

		StrictErrors: false,

		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer {
				Timeout: 2 * time.Second,
			}

			return d.DialContext(ctx, "udp", dnsServerAddr)
		}
	}
}
