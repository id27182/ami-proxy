package asmi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type responseJson struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn string `json:"expires_in"`
	ExpiresOn string `json:"expires_on"`
	NotBefore string `json:"not_before"`
	Resource string `json:"resource"`
	TokenType string `json:"token_type"`
}

func GetAmiToken(resource string) (string, error)  {
	// validate parameters
	if resource == "" {
		return "", fmt.Errorf("empty reqource is not allowed")
	}

	// Create HTTP request for a managed services for Azure resources token to access Azure Resource Manager
	var msi_endpoint *url.URL
	msi_endpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01")
	if err != nil {
		return "", fmt.Errorf("unable to create url for request. Original error: %s", err)
	}

	msi_parameters := url.Values{}
	msi_parameters.Add("resource", resource)
	msi_parameters.Add("api-version", "2018-02-01")
	msi_endpoint.RawQuery = msi_parameters.Encode()
	req, err := http.NewRequest("GET", msi_endpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("unable to create http request. Original error: %s", err)
	}
	req.Header.Add("Metadata", "true")

	// Call managed services for Azure resources token endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		return "", fmt.Errorf("unable to call token endpoint. Original error: %s", err)
	}

	// Pull out response body
	responseBytes,err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("unable to read response body. Original error: %s", err)
	}

	// Unmarshall response body into struct
	var r responseJson
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshall response. Original error: %s", err)
	}

	return r.AccessToken, nil
}
