package server

import (
	"fmt"
	"github.com/id27182/ami-proxy/pkg/asmi"
	"io"
	"log"
	"net/http"
)

type AmiProxy struct {
	DestHost string
	DestPort string
	DestProtocol string
	DestResource string
}

func (p *AmiProxy)  ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	transport.DisableCompression = true
	client.Transport = transport

	// rewrite dest protocol, host, and port
	destUri := fmt.Sprintf("%s://%s:%s%s", p.DestProtocol, p.DestHost, p.DestPort, r.URL.Path)
	log.Printf(destUri)

	req, err = http.NewRequest(r.Method, destUri, r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}

	amiToken, err := asmi.GetAmiToken(p.DestResource)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("token: %s", amiToken)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", amiToken))


	resp, err = client.Do(req)
	if err != nil {
		log.Printf(err.Error())
	}
	r.Body.Close()

	// combined for GET/POST
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	//conn := &HttpConnection{r, resp}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

	// Print HTTP conn
	//connChannel <-  conn
}


