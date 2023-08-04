package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/laof/lite-speed-test/config"
	_ "github.com/laof/lite-speed-test/config"
	"github.com/laof/lite-speed-test/dns"
	"github.com/laof/lite-speed-test/outbound"
	"github.com/laof/lite-speed-test/proxy"
	"github.com/laof/lite-speed-test/request"
	"github.com/laof/lite-speed-test/transport/resolver"
	"github.com/laof/lite-speed-test/tunnel"
	"github.com/laof/lite-speed-test/tunnel/adapter"
	"github.com/laof/lite-speed-test/tunnel/http"
	"github.com/laof/lite-speed-test/tunnel/socks"
	"github.com/laof/lite-speed-test/utils"
)

func StartInstance(c Config) (*proxy.Proxy, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "LocalHost", c.LocalHost)
	ctx = context.WithValue(ctx, "LocalPort", c.LocalPort)
	adapterServer, err := adapter.NewServer(ctx, nil)
	if err != nil {
		return nil, err
	}
	httpServer, err := http.NewServer(ctx, adapterServer)
	if err != nil {
		return nil, err
	}
	socksServer, err := socks.NewServer(ctx, adapterServer)
	if err != nil {
		return nil, err
	}
	sources := []tunnel.Server{httpServer, socksServer}
	sink, err := createSink(ctx, c.Link)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func(link string) {
		if c.Ping < 1 {
			return
		}
		if cfg, err := config.Link2Config(c.Link); err == nil {
			opt := request.PingOption{
				Attempts: c.Ping,
				TimeOut:  1200 * time.Millisecond,
			}
			info := fmt.Sprintf("%s %s", cfg.Remarks, net.JoinHostPort(cfg.Server, strconv.Itoa(cfg.Port)))
			if elapse, err := request.PingLinkInternal(link, opt); err == nil {
				info = fmt.Sprintf("%s \033[32m%dms\033[0m", info, elapse)
			} else {
				info = fmt.Sprintf("\033[31m%s\033[0m", err.Error())
			}
			log.Print(info)
		}
	}(c.Link)
	setDefaultResolver()
	p := proxy.NewProxy(ctx, cancel, sources, sink)
	return p, nil
}

func createSink(ctx context.Context, link string) (tunnel.Client, error) {
	var d outbound.Dialer
	matches, err := utils.CheckLink(link)
	if err != nil {
		return nil, err
	}
	creator, err := outbound.GetDialerCreator(matches[1])
	if err != nil {
		return nil, err
	}
	d, err = creator(link)
	if err != nil {
		return nil, err
	}
	if d != nil {
		return proxy.NewClient(ctx, d), nil
	}

	return nil, errors.New("not supported link")
}

func setDefaultResolver() {
	resolver.DefaultResolver = dns.DefaultResolver()
}
