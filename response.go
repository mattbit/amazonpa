package amazonpa

import "encoding/xml"

// Response describes the generic API Response
type Response struct {
	OperationRequest struct {
		RequestID             string     `xml:"RequestId"`
		Arguments             []Argument `xml:"Arguments>Argument"`
		RequestProcessingTime float64
	}
}

// Argument todo
type Argument struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

// Image todo
type Image struct {
	XMLName xml.Name `xml:"MediumImage"`
	URL     string
	Height  uint16
	Width   uint16
}

// Price describes the product price as
// Amount of cents in CurrencyCode
type Price struct {
	Amount         uint
	CurrencyCode   string
	FormattedPrice string
}

// Item represents a product returned by the API
type Item struct {
	ASIN           string
	URL            string
	DetailPageURL  string
	ItemAttributes ItemAttributes
	OfferSummary   OfferSummary
	Offers         Offers
}

// ItemAttributes response group
type ItemAttributes struct {
	Brand        string
	Title        string
	Feature      string
	Manufacturer string
	Model        string
	ProductGroup string
	Publisher    string
	ReleaseDate  string
	Studio       string
	Warranty     string
}

// Offer response attribute
type Offer struct {
	Condition       string `xml:"OfferAttributes>Condition"`
	ID              string `xml:"OfferListing>OfferListingId"`
	Price           Price  `xml:"OfferListing>Price"`
	PercentageSaved uint   `xml:"OfferListing>PercentageSaved"`
	Availability    string `xml:"OfferListing>Availability"`
}

// Offers response group
type Offers struct {
	TotalOffers     int
	TotalOfferPages int
	MoreOffersURL   string  `xml:"MoreOffersUrl"`
	Offers          []Offer `xml:"Offer"`
}

// OfferSummary response group
type OfferSummary struct {
	LowestNewPrice   Price
	LowerUsedPrice   Price
	TotalNew         int
	TotalUsed        int
	TotalCollectible int
	TotalRefurbished int
}

// ItemLookupRequest is the confirmation of a ItemLookup request
type ItemLookupRequest struct {
	IDType         string   `xml:"IdType"`
	ItemID         string   `xml:"ItemId"`
	ResponseGroups []string `xml:"ResponseGroup"`
	VariationPage  string
}

// ItemLookupResponse describes the API response for the ItemLookup operation
type ItemLookupResponse struct {
	Response

	Items struct {
		Request struct {
			IsValid           bool
			ItemLookupRequest ItemLookupRequest
		}
		Item Item
	}
}
