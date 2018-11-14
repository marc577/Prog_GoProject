package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirIfNotExist(t *testing.T) {

	assert.Equal(t, createDirIfNotExist("../../log"), true, "Folder Creation should be possible")

}

func TestCloseLogfile(t *testing.T) {

	file := createLogfileIfNotExist("test")
	assert.Equal(t, closeLogFile(file), true, "LogFile Closing should be possible")

}
