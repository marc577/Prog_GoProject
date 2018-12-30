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
	assert.False(t, createDirIfNotExist("!@#$%^&*()_"))
}

func TestCreateUserJSONIfNotExist(t *testing.T) {
	assert.True(t, createUserJSONIfNotExist("../../users.json"))
	os.Remove("../../users.json")
	assert.False(t, createDirIfNotExist("!@#$%^&*()_.json"))
}

func TestStartup(t *testing.T) {
	//startup("../../../log", "../../../storage/users.json", "../../../storage/tickets/", 8443, "../../../keys/server.crt", "../../../keys/server.key", "../../../html")
}
