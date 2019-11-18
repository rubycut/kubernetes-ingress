package main

import (
	parser "github.com/haproxytech/config-parser/v2"
	tcptype "github.com/haproxytech/config-parser/v2/parsers/tcp/actions"
	"github.com/haproxytech/config-parser/v2/types"
)

func (c *HAProxyController) HAProxyTCPInitialize() {
	config, errParser := c.ActiveConfiguration()
	PanicErr(errParser)

	for _, frontend := range []string{FrontendHTTP, FrontendHTTPS} {
		errParser = config.Set(parser.Frontends, frontend, "mode", types.StringC{Value: "tcp"})
		PanicErr(errParser)
		errParser = config.Insert(parser.Frontends, frontend, "tcp-request", &tcptype.Content{ //must be a pointer!
			Action:   []string{"accept"},
			Cond:     "if",
			CondTest: "{ req_ssl_hello_type 1 }",
			Comment:  "Added on func HAProxyController->HAProxyInitialize",
		}, 0)
		PanicErr(errParser)

		errParser = config.Insert(parser.Frontends, frontend, "tcp-request", &tcptype.InspectDelay{ //must be a pointer!
			Timeout: "5s",
			Comment: "Added on func HAProxyController->HAProxyInitialize",
		}, 0)
		PanicErr(errParser)
	}

	err := config.Delete(parser.Defaults, parser.DefaultSectionName, "option httplog")
	PanicErr(err)
	err = config.Set(parser.Defaults, parser.DefaultSectionName, "option tcplog", types.SimpleOption{
		Comment: "Added on func HAProxyController->HAProxyInitialize",
	})
	PanicErr(err)
}
