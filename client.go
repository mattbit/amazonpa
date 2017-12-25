package amazonpa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// ItemLookupQuery describes the allowed parameters for a ItemLookup request
type ItemLookupQuery struct {
	Condition             string
	IDType                string
	IncludeReviewsSummary string
	ItemIDs               []string
	MerchantID            string
	RelatedItemPage       string
	RelationshipType      string
	SearchIndex           string
	TruncateReviewsAt     string
	VariationPage         string
	ResponseGroups        []string
}

// ItemSearchQuery describes the allowed parameters for a ItemSearch request
type ItemSearchQuery struct {
	Actor                 string
	Artist                string
	AudienceRatings       []string
	Author                string
	Availability          string
	Brand                 string
	BrowseNode            string
	Composer              string
	Condition             string
	Conductor             string
	Director              string
	IncludeReviewsSummary string
	ItemPage              string
	Keywords              string
	Manufacturer          string
	MaximumPrice          string
	MerchantID            string
	MinimumPrice          string
	MinPercentageOff      string
	Orchestra             string
	Power                 string
	Publisher             string
	RelatedItemPage       string
	RelationshipType      string
	SearchIndex           string
	Sort                  string
	Title                 string
	TruncateReviewsAt     string
	VariationPage         string
	ResponseGroups        []string
}

type BrowseNodeLookupQuery struct {
	BrowseNodeID   string
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
func (client Client) ProcessRequest(request *Request, timeout time.Duration) ([]byte, error) {

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

	httpClient := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, errors.New("amazonpa: error on building request")
	}

	httpResponse, err = httpClient.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil, errors.New(fmt.Sprintf("amazonpa: error on request: %s - error: %s", requestURL, err.Error()))
	}
	contents, err = ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()

	if err != nil {
		return nil, errors.New("amazonpa: error while reading the server response")
	}

	return contents, nil
}

// ItemLookupWithTimeout performs an ItemLookup request with timeout specified
func (client Client) ItemLookupWithTimeout(query ItemLookupQuery, timeout time.Duration) (*ItemLookupResponse, error) {

	request := client.NewRequest("ItemLookup")

	request.SetParameter("Condition", query.Condition)
	request.SetParameter("IdType", query.IDType)
	request.SetParameter("IncludeReviewsSummary", query.IncludeReviewsSummary)
	request.SetParameter("ItemId", strings.Join(query.ItemIDs, ","))
	request.SetParameter("MerchantId", query.MerchantID)
	request.SetParameter("RelatedItemPage", query.RelatedItemPage)
	request.SetParameter("RelationshipType", query.RelationshipType)
	request.SetParameter("SearchIndex", query.SearchIndex)
	request.SetParameter("TruncateReviewsAt", query.TruncateReviewsAt)
	request.SetParameter("VariationPage", query.VariationPage)
	request.SetParameter("ResponseGroup", strings.Join(query.ResponseGroups, ","))

	xmlData, err := client.ProcessRequest(request, timeout)

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

// ItemLookup performs an ItemLookup request using DefaultTimeout
func (client Client) ItemLookup(query ItemLookupQuery) (*ItemLookupResponse, error) {
	return client.ItemLookupWithTimeout(query, client.config.DefaultTimeout)
}

// ItemLookup performs an ItemLookup request using DefaultTimeout
func (client Client) ItemSearch(query ItemSearchQuery) (*ItemSearchResponse, error) {
	return client.ItemSearchWithTimeout(query, client.config.DefaultTimeout)
}

// ItemLookup performs an ItemLookup request using specified Timeout
func (client Client) ItemSearchWithTimeout(query ItemSearchQuery, timeout time.Duration) (*ItemSearchResponse, error) {

	request := client.NewRequest("ItemSearch")

	request.SetParameter("Actor", query.Actor)
	request.SetParameter("Artist", query.Artist)
	request.SetParameter("AudienceRating", strings.Join(query.AudienceRatings, ","))
	request.SetParameter("Author", query.Author)
	request.SetParameter("Availability", query.Availability)
	request.SetParameter("Brand", query.Brand)
	request.SetParameter("BrowseNode", query.BrowseNode)
	request.SetParameter("Composer", query.Composer)
	request.SetParameter("Condition", query.Condition)
	request.SetParameter("Conductor", query.Conductor)
	request.SetParameter("Director", query.Director)
	request.SetParameter("IncludeReviewsSummary", query.IncludeReviewsSummary)
	request.SetParameter("ItemPage", query.ItemPage)
	request.SetParameter("Keywords", query.Keywords)
	request.SetParameter("Manufacturer", query.Manufacturer)
	request.SetParameter("MaximumPrice", query.MaximumPrice)
	request.SetParameter("MerchantId", query.MerchantID)
	request.SetParameter("MinimumPrice", query.MinimumPrice)
	request.SetParameter("MinPercentageOff", query.MinPercentageOff)
	request.SetParameter("Orchestra", query.Orchestra)
	request.SetParameter("Power", query.Power)
	request.SetParameter("Publisher", query.Publisher)
	request.SetParameter("RelatedItemPage", query.RelatedItemPage)
	request.SetParameter("RelationshipType", query.RelationshipType)
	request.SetParameter("SearchIndex", query.SearchIndex)
	request.SetParameter("Sort", query.Sort)
	request.SetParameter("Title", query.Title)
	request.SetParameter("TruncateReviewsAt", query.TruncateReviewsAt)
	request.SetParameter("VariationPage", query.VariationPage)
	request.SetParameter("ResponseGroup", strings.Join(query.ResponseGroups, ","))

	xmlData, err := client.ProcessRequest(request, timeout)

	if err != nil {
		return nil, err
	}

	var response ItemSearchResponse
	xml.Unmarshal(xmlData, &response)

	if response.Items.Request.IsValid != true {
		return &response, errors.New("amazonpa: request is invalid")
	}

	return &response, nil
}

// ItemLookup performs an ItemLookup request with default timeout
func (client Client) BrowseNodeLookup(query BrowseNodeLookupQuery) (*BrowseNodeLookupResponse, error) {
	return client.BrowseNodeLookupWithTimeout(query, client.config.DefaultTimeout)
}

// ItemLookup performs an ItemLookup request using specified timeout
func (client Client) BrowseNodeLookupWithTimeout(query BrowseNodeLookupQuery, timeout time.Duration) (*BrowseNodeLookupResponse, error) {

	request := client.NewRequest("BrowseNodeLookup")

	request.SetParameter("BrowseNodeId", query.BrowseNodeID)
	request.SetParameter("ResponseGroup", strings.Join(query.ResponseGroups, ","))

	xmlData, err := client.ProcessRequest(request, timeout)

	if err != nil {
		return nil, err
	}

	var response BrowseNodeLookupResponse
	xml.Unmarshal(xmlData, &response)

	if response.BrowseNodes.Request.IsValid != true {
		return &response, errors.New("amazonpa: request is invalid")
	}

	return &response, nil
}
