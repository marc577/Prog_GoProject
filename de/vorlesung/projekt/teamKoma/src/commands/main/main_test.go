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
	os.Mkdir("../../test/", 0400)
	assert.False(t, createDirIfNotExist("../../test/test"))
	os.RemoveAll("../../test/")
}
