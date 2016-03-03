package models

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var JSON = []byte(`{
  "dcdr": {
  	"info" : {
  		"current_sha" : "asdfasd"
  	},
    "features": {
      "default":{
        "test": 0.1,
        "bool": true
      },
      "cc": {
      	"test": false
      }
    }
  }
}`)

func TestJSON(t *testing.T) {
	var d DcdrMap

	err := json.Unmarshal(JSON, &d)

	assert.NoError(t, err)

	fmt.Println(err)
	fmt.Printf("%+v", d.Dcdr)

	bts, _ := json.MarshalIndent(d, "", "  ")
	fmt.Println(string(bts[:]))
}
