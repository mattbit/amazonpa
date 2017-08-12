package amazonpa

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func assertEqualStr(t *testing.T, v1, v2 string, msg ...string) {
	if v1 != v2 {
		t.Error(msg)
	}
}

func assertEqualInt(t *testing.T, v1, v2 int, msg ...string) {
	if v1 != v2 {
		t.Error(msg)
	}
}

func assertEqualFloat(t *testing.T, v1, v2 float64, msg ...string) {
	if v1 != v2 {
		t.Error(msg)
	}
}

func assertEqualBool(t *testing.T, v1, v2 bool, msg ...string) {
	if v1 != v2 {
		t.Error(msg)
	}
}

func TestParseItemLookupResponse(t *testing.T) {
	responseData, err := ioutil.ReadFile("testdata/itemlookup_response.xml")
	if err != nil {
		t.Error(err)
	}
	var response ItemLookupResponse
	xml.Unmarshal([]byte(responseData), &response)

	// OperationRequest
	assertEqualStr(t, response.OperationRequest.RequestID, "54b32dd2-f525-4987-9295-b1f2cda8d5f5", "Bad RequestID")
	assertEqualStr(t, response.OperationRequest.Arguments[0].Name, "ItemId", "Bad Argument Name")
	assertEqualStr(t, response.OperationRequest.Arguments[0].Value, "B003TGG2EA", "Bad Argument Value")
	assertEqualFloat(t, response.OperationRequest.RequestProcessingTime, 0.0837001980000000, "Bad RequestProcessingTime")

	// Items/Request
	assertEqualBool(t, response.Items.Request.IsValid, true, "Bad IsValid")
	assertEqualStr(t, response.Items.Request.ItemLookupRequest.IDType, "ASIN", "Bad IDType")
	assertEqualStr(t, response.Items.Request.ItemLookupRequest.ItemID, "B003TGG2EA", "Bad ItemID")
	assertEqualStr(t, response.Items.Request.ItemLookupRequest.ResponseGroup, "Large", "Bad ResponseGroup")
	assertEqualStr(t, response.Items.Request.ItemLookupRequest.VariationPage, "All", "Bad VariationPage")

	//Items/Item
	assertEqualStr(t, response.Items.Item.ASIN, "B003TGG2EA", "Bad ASIN")
	assertEqualStr(t, response.Items.Item.DetailPageURL, "https://www.amazon.it/Grohe-32843000-Cosmopolitan-Miscelatore-Monocomando/dp/B003TGG2EA%3FSubscriptionId%3DAKIAIZL74FSKHXDX66WQ%26tag%3Dgoldbot-21%26linkCode%3Dxm2%26camp%3D2025%26creative%3D165953%26creativeASIN%3DB003TGG2EA", "Bad DetailPageURL")
	assertEqualInt(t, response.Items.Item.SalesRank, 214, "Bad SalesRank")

	assertEqualStr(t, response.Items.Item.SmallImage.URL, "http://ecx.images-amazon.com/images/I/3143fcCtWnL._SL75_.jpg", "Bad SmallImage/URL")
	assertEqualInt(t, int(response.Items.Item.SmallImage.Height), 75, "Bad SmallImage/Height")
	assertEqualInt(t, int(response.Items.Item.SmallImage.Width), 52, "Bad SmallImage/Width")

	assertEqualStr(t, response.Items.Item.MediumImage.URL, "http://ecx.images-amazon.com/images/I/3143fcCtWnL._SL160_.jpg", "Bad SmallImage/URL")
	assertEqualInt(t, int(response.Items.Item.MediumImage.Height), 160, "Bad SmallImage/Height")
	assertEqualInt(t, int(response.Items.Item.MediumImage.Width), 110, "Bad SmallImage/Width")

	assertEqualStr(t, response.Items.Item.LargeImage.URL, "http://ecx.images-amazon.com/images/I/3143fcCtWnL.jpg", "Bad SmallImage/URL")
	assertEqualInt(t, int(response.Items.Item.LargeImage.Height), 500, "Bad SmallImage/Height")
	assertEqualInt(t, int(response.Items.Item.LargeImage.Width), 344, "Bad SmallImage/Width")

	assertEqualStr(t, response.Items.Item.ImageSets.ImageSet[0].Category, "primary", "Bad ImageSet category")
	assertEqualInt(t, int(response.Items.Item.ImageSets.ImageSet[0].SwatchImage.Height), 30, "Bad SwatchImage")

	// Items/Item/ItemAttributes
	assertEqualStr(t, response.Items.Item.ItemAttributes.Binding, "Tools & Home Improvement", "Bad Binding")
	assertEqualStr(t, response.Items.Item.ItemAttributes.Brand, "Grohe", "Bad Binding")
	assertEqualStr(t, response.Items.Item.ItemAttributes.Color, "Cromo", "Bad Color")
	assertEqualStr(t, response.Items.Item.ItemAttributes.Creator, "", "Bad Creator")
	assertEqualStr(t, response.Items.Item.ItemAttributes.EAN, "4005176874840", "Bad EAN")

	assertEqualInt(t, int(response.Items.Item.ItemAttributes.ListPrice.Amount), 18500, "Bad Price/Amount")
	assertEqualStr(t, response.Items.Item.ItemAttributes.ListPrice.CurrencyCode, "EUR", "Bad Price/Currency")
	assertEqualStr(t, response.Items.Item.ItemAttributes.ListPrice.FormattedPrice, "EUR 185,00", "Bad Price/Formatted Price")

	assertEqualStr(t, response.Items.Item.BrowseNodes.BrowseNode[0].BrowseNodeID, "3120323031", "Bad BrowseNode/BrowseNodeID")
	assertEqualStr(t, response.Items.Item.BrowseNodes.BrowseNode[0].Name, "Rubinetti per lavelli da cucina", "Bad BrowseNode/Name")

	assertEqualStr(t, response.Items.Item.BrowseNodes.BrowseNode[0].Ancestors.BrowseNode[0].BrowseNodeID, "3119756031", "Bad Ancestors/BrowseNodeID")
	assertEqualStr(t, response.Items.Item.BrowseNodes.BrowseNode[0].Ancestors.BrowseNode[0].Name, "Rubinetti da cucina", "Bad Ancestors/Name")
}
