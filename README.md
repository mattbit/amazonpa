[![Build Status](https://travis-ci.org/mattbit/amazonpa.svg?branch=master)](https://travis-ci.org/mattbit/amazonpa)

# Amazon PA

A Go lang library to interact with Amazon Product Advertising API.

The library does not cover all the functionality of the PA API, but feel free to create new pull requests and I will try to merge them quickly!

## Example

```go
package main

import (
	"fmt"

	"github.com/mattbit/amazonpa"
)

func main() {
	cfg := amazonpa.Config{
		AccessKey:    "YOUR_KEY",
		AccessSecret: "YOUR_SECRET",
		AssociateTag: "YOUR_TAG",
		Region:       "YOUR_REGION",
		Secure:       true,
	}
	client := amazonpa.NewClient(cfg)

	query := amazonpa.ItemSearchQuery{
		SearchIndex:    "All",
		Keywords:       "mouse",
		ResponseGroups: []string{"Large"},
	}

	response, err := client.ItemSearch(query)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Found %d items\n", len(response.Items.Items))
	}
}
```
