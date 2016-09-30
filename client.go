package amazonpa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ItemLookupQuery describes the allowed parameters for a ItemLookup request
type ItemLookupQuery struct {
	ItemIDs        []string
	IDType         string
	MerchantID     string
	ResponseGroups []string
}

// Client provides the functions to interact with the API
type Client struct {
	config Config
}

// NewClient returns a new Client
func NewClient(config Config) *Client {
	c := Client{config}

	return &c
}

// NewRequest returns a request with basic parameters
func (client Client) NewRequest(operation string) *Request {

	request := Request{}

	if client.config.Secure {
		request.scheme = "https"
	} else {
		request.scheme = "http"
	}

	request.endpoint = Endpoints[client.config.Region]
	request.endpointURI = EndpointURI

	request.parameters = map[string]string{
		"Service":        "AWSECommerceService",
		"AWSAccessKeyId": client.config.AccessKey,
		"AssociateTag":   client.config.AssociateTag,
		"Version":        "2013-08-01",
		"Operation":      operation,
		"Timestamp":      time.Now().Format(time.RFC3339),
	}

	return &request
}

// SignRequest produces the signature for the given query string
func (client Client) SignRequest(request *Request) {
	signable := fmt.Sprintf("GET\n%s\n%s\n%s", request.endpoint, request.endpointURI, request.QueryString())

	hasher := hmac.New(sha256.New, []byte(client.config.AccessSecret))
	hasher.Write([]byte(signable))

	request.signature = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// ProcessRequest takes a request and queries the API
func (client Client) ProcessRequest(request *Request) ([]byte, error) {

	// Sign the request
	client.SignRequest(request)

	// Execute the HTTP request
	var httpResponse *http.Response
	var err error
	var contents []byte

	requestURL, err := request.SignedURL()

	if err != nil {
		return nil, errors.New("amazonpa: cannot get the signed request URL")
	}

	httpResponse, err = http.Get(requestURL)

	if err != nil {
		return nil, errors.New("amazonpa: error processing the http request")
	}

	contents, err = ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()

	if err != nil {
		return nil, errors.New("amazonpa: error while reading the server response")
	}

	return contents, nil
}

// ItemLookup performs an ItemLookup request
func (client Client) ItemLookup(query ItemLookupQuery) (*ItemLookupResponse, error) {

	request := client.NewRequest("ItemLookup")

	request.SetParameter("ItemId", strings.Join(query.ItemIDs, ","))
	request.SetParameter("ResponseGroup", strings.Join(query.ResponseGroups, ","))
	request.SetParameter("IdType", query.IDType)
	request.SetParameter("SearchIndex", "All")
	request.SetParameter("MerchantId", query.MerchantID)

	xmlData, err := client.ProcessRequest(request)

	if err != nil {
		return nil, err
	}

	var response ItemLookupResponse
	xml.Unmarshal(xmlData, &response)

	if response.Items.Request.IsValid != true {
		return &response, errors.New("amazonpa: request is invalid")
	}

	return &response, nil
}

// ItemLookup performs an ItemLookup request
func (client Client) ItemLookupByKeyword(keywords []string) (*ItemLookupResponse, error) {

	request := client.NewRequest("ItemSearch")

	request.SetParameter("ResponseGroup", "Large")
	request.SetParameter("SearchIndex", "All")
	request.SetParameter("Keywords", strings.ToLower((strings.Join(keywords, "_"))))

	xmlData, err := client.ProcessRequest(request)

	if err != nil {
		return nil, err
	}

	var response ItemLookupResponse
	xml.Unmarshal(xmlData, &response)

	if response.Items.Request.IsValid != true {
		return &response, errors.New("amazonpa: request is invalid")
	}

	return &response, nil
}
