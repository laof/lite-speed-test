package http

import (
	"context"

	"github.com/laof/lite-speed-test/tunnel"
)

const Name = "HTTP"

type Tunnel struct{}

func (t *Tunnel) Name() string {
	return Name
}

func (t *Tunnel) NewServer(ctx context.Context, server tunnel.Server) (tunnel.Server, error) {
	return NewServer(ctx, server)
}
