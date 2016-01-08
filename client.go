package amazonpa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client provides the functions to interact with the API
type Client struct {
	accessKeyID  string
	secretKey    string
	associateTag string
	region       string
	secure       bool
}

// NewClient returns a new Client
func NewClient(accessKeyID string, secretKey string, associateTag string, region string, secure bool) *Client {
	c := Client{accessKeyID, secretKey, associateTag, region, secure}

	return &c
}

// NewRequest returns a request with basic parameters
func (client Client) NewRequest(operation string) *Request {

	request := Request{}

	if client.secure {
		request.scheme = "https"
	} else {
		request.scheme = "http"
	}

	request.endpoint = Endpoints[client.region]
	request.endpointURI = EndpointURI

	request.parameters = map[string]string{
		"Service":        "AWSECommerceService",
		"AWSAccessKeyId": client.accessKeyID,
		"AssociateTag":   client.associateTag,
		"Version":        "2013-08-01",
		"Operation":      operation,
		"Timestamp":      time.Now().Format(time.RFC3339),
	}

	return &request
}

// SignRequest produces the signature for the given query string
func (client Client) SignRequest(request *Request) {
	signable := fmt.Sprintf("GET\n%s\n%s\n%s", request.endpoint, request.endpointURI, request.QueryString())

	hasher := hmac.New(sha256.New, []byte(client.secretKey))
	hasher.Write([]byte(signable))

	request.signature = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// ProcessRequest takes a request and queries the API
func (client Client) ProcessRequest(request *Request) (string, error) {

	// Sign the request
	client.SignRequest(request)

	// Execute the HTTP request
	var httpResponse *http.Response
	var err error
	var contents []byte

	requestURL, err := request.SignedURL()

	if err != nil {
		return "", errors.New("amazonpa: cannot get the signed request URL")
	}

	httpResponse, err = http.Get(requestURL)

	if err != nil {
		return "", errors.New("amazonpa: error processing the http request")
	}

	contents, err = ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()

	if err != nil {
		return "", errors.New("amazonpa: response error")
	}

	return string(contents[:]), nil
}

// ItemLookup performs an ItemLookup request
func (client Client) ItemLookup(itemIDs []string, responseGroups []string, itemType string) {

	request := client.NewRequest("ItemLookup")

	request.SetParameter("ItemId", strings.Join(itemIDs, ","))
	request.SetParameter("ResponseGroup", strings.Join(responseGroups, ","))
	// request.SetParameter("ItemType", itemType)

	client.ProcessRequest(request)
}
