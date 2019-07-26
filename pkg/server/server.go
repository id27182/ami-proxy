package server

import (
	"fmt"
	"github.com/id27182/ami-proxy/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Serve() error {
	conf, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("unable to get config. Original error: %s", err)
	}

	proxyConfig, err := conf.Proxy()
	if err != nil {
		return fmt.Errorf("unable to get proxy config. Original error: %s")
	}

	amiProxy := AmiProxy{
		DestProtocol: proxyConfig.DestProtocol,
		DestHost: proxyConfig.DestHost,
		DestPort: proxyConfig.DestPort,
		DestResource: proxyConfig.DestResource,
	}

	log.Infof("starting proxy http listener on port: %s. Dest host: %s, dest port: %s, dest protocol: %s, dest resource: %s", proxyConfig.BindPort, proxyConfig.DestHost, proxyConfig.DestPort, proxyConfig.DestProtocol, proxyConfig.DestResource)
	err = http.ListenAndServe(fmt.Sprintf(":%s", proxyConfig.BindPort), &amiProxy)
	if err != nil {
		return fmt.Errorf("unable to start proxy http listener. Original error: %s", err)
	}

	return nil
}
