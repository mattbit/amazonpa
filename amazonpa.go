package amazonpa

// Endpoints are the Amazon API endpoints by region
var Endpoints = map[string]string{
	"CA": "ecs.amazonaws.ca",
	"CN": "webservices.amazon.cn",
	"DE": "ecs.amazonaws.de",
	"ES": "webservices.amazon.es",
	"FR": "ecs.amazonaws.fr",
	"IT": "webservices.amazon.it",
	"JP": "ecs.amazonaws.jp",
	"UK": "ecs.amazonaws.co.uk",
	"US": "webservices.amazon.com",
}

// EndpointURI is the fixed request URI of the API
const EndpointURI = "/onca/xml"
