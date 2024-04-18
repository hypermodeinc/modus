/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

func adjustResponse(r []byte) []byte {
	/*
		Without this function, errors would appear like this:
		{
		  "errors": [
		    {
		      "message": "Failed to fetch from Subgraph at Path 'query'.",
		      "extensions": {
		        "errors": [
		          {
		            "message": "This is an error message!",
		            "path": ["myFunction"],
		            "extensions": {"level": "error"}
		          }
		        ]
		      }
		    }
		  ]
		}

		This function unwraps the outer error so that it looks like this:
		{
		  "errors": [
		    {
		      "message": "This is an error message!",
		      "path": ["myFunction"],
		      "extensions": {"level": "error"}
		    }
		  ]
		}

	*/

	errors, _, _, err := jsonparser.Get(r, "errors")
	if err != nil {
		return r
	}

	i := 0
	jsonparser.ArrayEach(errors, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if s, err := jsonparser.GetString(value, "message"); err == nil {
			if strings.HasPrefix(s, "Failed to fetch from Subgraph") {
				errors, _, _, err := jsonparser.Get(value, "extensions", "errors")
				if err != nil {
					return
				}
				r, _ = jsonparser.Set(r, errors[1:len(errors)-1], "errors", fmt.Sprintf("[%d]", i))
			}
		}
		i++
	})

	return r
}
