package main

import (
	parser "github.com/haproxytech/config-parser/v2"
	tcptype "github.com/haproxytech/config-parser/v2/parsers/tcp/actions"
	"github.com/haproxytech/config-parser/v2/types"
)

func (c *HAProxyController) handleSSLPassthrough() (needsReload bool, err error) {
	config, errParser := c.ActiveConfiguration()
	PanicErr(errParser)
	annSSLPassthrough, _ := GetValueFromAnnotations("ssl-passthrough") // ADD c.cfg.ConfigMap.Annotations
	if annSSLPassthrough.Status == EMPTY {
		return false, nil
	}
	if annSSLPassthrough.Value == ENABLED {
		for _, frontend := range []string{FrontendHTTP, FrontendHTTPS} {
			errParser = config.Set(parser.Frontends, frontend, "mode", types.StringC{Value: "tcp"})
			PanicErr(errParser)

			errParser = config.Insert(parser.Frontends, frontend, "tcp-request", &tcptype.InspectDelay{ //must be a pointer!
				Timeout: "5s",
				Comment: "Added on func HAProxyController->HAProxyInitialize",
			}, 0)
			PanicErr(errParser)

			errParser = config.Set(parser.Frontends, frontend, "option tcplog", types.SimpleOption{
				Comment: "Added on func HAProxyController->HAProxyInitialize",
			})
			PanicErr(errParser)
		}
		errParser = config.Insert(parser.Frontends, FrontendHTTPS, "tcp-request", &tcptype.Content{ //must be a pointer!
			Action:   []string{"accept"},
			Cond:     "if",
			CondTest: "{ req_ssl_hello_type 1 }",
			Comment:  "Added on func HAProxyController->HAProxyInitialize",
		}, 0)
		PanicErr(errParser)
	} /*else {
		no else for now we either have it enabled and added or empty
	}*/
	return true, nil
}
