package amazonpa

import (
	"testing"
	"time"
)

func newTestClient() *Client {
	config := Config{
		"AKIAIOSFODNN7EXAMPLE",
		"1234567890",
		"mytag-20",
		"US",
		false,
	}

	return NewClient(config)
}

func TestRequestHasDefaultParameters(t *testing.T) {
	client := newTestClient()
	startTime := time.Now()
	request := client.NewRequest("ItemLookup")

	checkParameters := map[string]string{
		"Service":        "AWSECommerceService",
		"AWSAccessKeyId": "AKIAIOSFODNN7EXAMPLE",
		"AssociateTag":   "mytag-20",
		"Version":        "2013-08-01",
	}

	for key, value := range checkParameters {
		if request.parameters[key] != value {
			t.Errorf("Request parameter %s is wrong", key)
		}
	}

	requestTime, _ := time.Parse(time.RFC3339, request.parameters["Timestamp"])

	if !(requestTime.Unix() >= startTime.Unix() && requestTime.Unix() <= time.Now().Unix()) {
		t.Error("Request parameter Timestamp is invalid")
	}

}

func TestRequestURLAndSignature(t *testing.T) {
	client := newTestClient()
	request := client.NewRequest("ItemLookup")

	request.SetParameter("ItemId", "0679722769")
	request.SetParameter("ResponseGroup", "Images,ItemAttributes,Offers,Reviews")
	request.SetParameter("Timestamp", "2014-08-18T12:00:00Z")

	if request.QueryString() != "AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&AssociateTag=mytag-20&ItemId=0679722769&Operation=ItemLookup&ResponseGroup=Images%2CItemAttributes%2COffers%2CReviews&Service=AWSECommerceService&Timestamp=2014-08-18T12%3A00%3A00Z&Version=2013-08-01" {
		t.Error("Request query string is wrong")
	}

	client.SignRequest(request)

	if request.signature != "j7bZM0LXZ9eXeZruTqWm2DIvDYVUU3wxPPpp+iXxzQc=" {
		t.Error("Request signature is wrong")
	}

	requestURL, _ := request.SignedURL()

	if requestURL != "http://webservices.amazon.com/onca/xml?AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&AssociateTag=mytag-20&ItemId=0679722769&Operation=ItemLookup&ResponseGroup=Images%2CItemAttributes%2COffers%2CReviews&Service=AWSECommerceService&Timestamp=2014-08-18T12%3A00%3A00Z&Version=2013-08-01&Signature=j7bZM0LXZ9eXeZruTqWm2DIvDYVUU3wxPPpp%2BiXxzQc%3D" {
		t.Error("Request signed URL is wrong")
	}
}
