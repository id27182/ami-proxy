package server

import (
	"fmt"
	"github.com/id27182/ami-proxy/pkg/asmi"
	"github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
	"io"
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

	// Get azure service managed identity token
	amiToken, err := asmi.GetAmiToken(p.DestResource)
	if err != nil {
		log.Errorf(err.Error())
	} else {
		log.Debugf("token: %s", amiToken)
	}
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

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

	// Print HTTP conn
	conn := &HttpConnection{r, resp}
	conn.Log()
}


type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

func (conn *HttpConnection) Log() {
	// Generate connection uuid
	id, err := uuid.NewV4()
	if err != nil {
		log.Warnf("unable to generate uuid for log record. Original error: %s", err)
	}

	log.WithField("id", id).Infof("Received a request: %v %v\n", conn.Request.Method, conn.Request.RequestURI)
	log.WithField("id", id).Debugln("Request headers:")
	for k, v := range conn.Request.Header {
		log.WithField("id", id).Debug(k, ":", v)
	}

	log.WithField("id", id).Infof("Response: HTTP/1.1 %v\n", conn.Response.Status)
	log.WithField("id", id).Debugln("Response headers:")
	for k, v := range conn.Response.Header {
		log.WithField("id", id).Debug(k, ":", v)
	}
	log.WithField("id", id).Debug(conn.Response.Body)
}