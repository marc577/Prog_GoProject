//Matrikelnummern:
//9188103
//1798794
//4717960
package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGenJSONData(t *testing.T) {
	assert.NotNil(t, genJSONData("test@test.com", "test", "test"))
}
func TestSendReq(t *testing.T) {
	jsonData := genJSONData("test@test.de", "test", "test")
	assert.Nil(t, sendReq("localhost", 8443, jsonData))
}
