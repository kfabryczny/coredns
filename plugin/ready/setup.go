package ready

import (
	"net"

	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("ready", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	addr, err := parse(c)
	if err != nil {
		return plugin.Error("ready", err)
	}

	r := new(addr)

	c.OnStartup(r.onStartup)
	c.OnStartup(r.wait)
	c.OnRestart(r.onRestart)
	c.OnFinalShutdown(r.onFinalShutdown)

	return nil
}

func parse(c *caddy.Controller) (string, error) {
	addr := ""
	for c.Next() {
		args := c.RemainingArgs()

		switch len(args) {
		case 0:
		case 1:
			addr = args[0]
			if _, _, e := net.SplitHostPort(addr); e != nil {
				return "", e
			}
		default:
			return "", c.ArgErr()
		}
	}
	return addr, nil
}
