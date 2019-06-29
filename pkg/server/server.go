package server

import (
	"fmt"
	"github.com/id27182/ami-proxy/pkg/config"
	"net/http"
)

//type HttpConnection struct {
//	Request  *http.Request
//	Response *http.Response
//}
//
//type HttpConnectionChannel chan *HttpConnection
//
//var connChannel = make(HttpConnectionChannel)
//
//func PrintHTTP(conn *HttpConnection) {
//	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
//	for k, v := range conn.Request.Header {
//		fmt.Println(k, ":", v)
//	}
//	fmt.Println("==============================")
//	fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
//	for k, v := range conn.Response.Header {
//		fmt.Println(k, ":", v)
//	}
//	fmt.Println(conn.Response.Body)
//	fmt.Println("==============================")
//}
//
//func HandleHTTP() {
//	for {
//		select {
//		case conn := <-connChannel:
//			PrintHTTP(conn)
//		}
//	}
//}

func Serve() error {
	conf, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("unable to get config. Original error: %s")
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

	err = http.ListenAndServe(fmt.Sprintf(":%s", proxyConfig.BindPort), &amiProxy)
	if err != nil {
		return fmt.Errorf("unable to start proxy http listener. Original error: %s", err)
	}

	return nil
}
