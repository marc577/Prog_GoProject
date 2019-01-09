//Matrikelnummern:
//9188103
//1798794
//4717960
package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirIfNotExist(t *testing.T) {
	assert.True(t, createDirIfNotExist("../../test/"))
	os.Remove("../../test/")
	//assert.False(t, createDirIfNotExist("!@#$%^&*()_"))
}

func TestCreateUserJSONIfNotExist(t *testing.T) {
	suc, exist := createUserJSONIfNotExist("../../users.json")
	assert.NotNil(t, suc)
	assert.NotNil(t, exist)
	os.Remove("../../users.json")
	suc1, _ := createUserJSONIfNotExist("!@#$%^&*()_.json")
	assert.NotNil(t, suc1)
}
