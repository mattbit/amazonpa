package amazonpa

// Request provides a generic API request
import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Request represents an API request
type Request struct {
	scheme      string
	endpoint    string
	endpointURI string
	signature   string
	parameters  map[string]string
}

// SetParameter adds a parameter to the request
func (request Request) SetParameter(key string, value string) {
	request.parameters[key] = value
}

// Parameters returns the request parameters
func (request Request) Parameters() map[string]string {
	return request.parameters
}

// QueryString returns the query string of the request
func (request Request) QueryString() string {
	var queryString string

	// Sort the parameters keys
	var keys []string
	for key := range request.parameters {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build the query string
	for _, key := range keys {
		if url.QueryEscape(request.parameters[key]) == "" {
			continue // Skip empty params
		}
		queryString += fmt.Sprintf("%s=%s&", url.QueryEscape(key), url.QueryEscape(request.parameters[key]))
	}

	// Remove last '&'
	return strings.TrimRight(queryString, "&")
}

// URL gives the complete URL of the request without the signature
func (request Request) URL() string {
	return fmt.Sprintf("%s://%s%s?%s", request.scheme, request.endpoint, request.endpointURI, request.QueryString())
}

// SignedURL returns the full signed request URL
func (request Request) SignedURL() (string, error) {
	if request.signature == "" {
		return "", errors.New("amazonpa: request is unsigned")
	}

	return fmt.Sprintf("%s&Signature=%s", request.URL(), url.QueryEscape(request.signature)), nil
}
