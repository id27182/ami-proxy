package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

type HttpConnectionChannel chan *HttpConnection

var connChannel = make(HttpConnectionChannel)

func PrintHTTP(conn *HttpConnection) {
	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
	for k, v := range conn.Request.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println("==============================")
	fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
	for k, v := range conn.Response.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println(conn.Response.Body)
	fmt.Println("==============================")
}

func HandleHTTP() {
	for {
		select {
		case conn := <-connChannel:
			PrintHTTP(conn)
		}
	}
}

type AmiProxy struct {
	DestHost string
	DestPort string
	DestProtocol string
}

func (p *AmiProxy)  ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	transport.DisableCompression = true
	client.Transport = transport

	destUri := fmt.Sprintf("%s://%s:%s%s", p.DestProtocol, p.DestHost, p.DestPort, r.URL.Path)
	log.Printf(destUri)

	req, err = http.NewRequest(r.Method, destUri, r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}

	amiToken, err := getAmiToken()
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

	conn := &HttpConnection{r, resp}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

	// Print HTTP conn
	connChannel <-  conn
}

func main() {
	go HandleHTTP()

	proxyRules := []ProxyRule{
		{
			SourcePort: "3333",

			DestHost: "ami-proxy.azurewebsites.net",
			DestPort: "443",
			DestProtocol: "https",
		},
	}

	for _, proxyRule := range proxyRules {
		amiProxy := AmiProxy{
			DestHost: proxyRule.DestHost,
			DestPort: proxyRule.DestPort,
			DestProtocol: proxyRule.DestProtocol,
		}

		err := http.ListenAndServe(fmt.Sprintf(":%s", proxyRule.SourcePort), &amiProxy)
		if err != nil {
			log.Fatalf("Listen and serve error: %s", err.Error())
		}
	}
}

type ProxyRule struct {
	SourcePort string

	DestHost string
	DestPort string
	DestProtocol string
}



type responseJson struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn string `json:"expires_in"`
	ExpiresOn string `json:"expires_on"`
	NotBefore string `json:"not_before"`
	Resource string `json:"resource"`
	TokenType string `json:"token_type"`
}

func getAmiToken() (string, error)  {
	// Create HTTP request for a managed services for Azure resources token to access Azure Resource Manager
	var msi_endpoint *url.URL
	msi_endpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01")
	if err != nil {
		fmt.Println("Error creating URL: ", err)
		return "", err
	}

	msi_parameters := url.Values{}
	//msi_parameters.Add("resource", "https://ami-proxy.azurewebsites.net")
	msi_parameters.Add("resource", "64d9b013-b84f-4c52-a59d-df4bd7f3a3d5")
	msi_parameters.Add("api-version", "2018-02-01")
	msi_endpoint.RawQuery = msi_parameters.Encode()
	req, err := http.NewRequest("GET", msi_endpoint.String(), nil)
	if err != nil {
		fmt.Println("Error creating HTTP request: ", err)
		return "", err
	}
	req.Header.Add("Metadata", "true")

	// Call managed services for Azure resources token endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println("Error calling token endpoint: ", err)
		return "", err
 	}

	// Pull out response body
	responseBytes,err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error reading response body : ", err)
		return "", err
	}

	// Unmarshall response body into struct
	var r responseJson
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		fmt.Println("Error unmarshalling the response:", err)
		return "", err
	}

	var r1 interface{}
	err = json.Unmarshal(responseBytes, &r1)
	fmt.Printf("%s", r1)

	fmt.Println("Response status:", resp.Status)
	log.Println("access_token: ", r.AccessToken)
	fmt.Println("refresh_token: ", r.RefreshToken)
	fmt.Println("expires_in: ", r.ExpiresIn)
	fmt.Println("expires_on: ", r.ExpiresOn)
	fmt.Println("not_before: ", r.NotBefore)
	fmt.Println("resource: ", r.Resource)
	fmt.Println("token_type: ", r.TokenType)

	return r.AccessToken, nil
}