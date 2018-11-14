package logging

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitQueues(t *testing.T) {
	assert.Equal(t, true, initQueues(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr), "Logging Output handles should be possible to renew")
}

func TestLogInit(t *testing.T) {
	assert.Equal(t, true, LogInit("../../log"), "explicit logging procedure could not be started")
}

func TestCreateDirIfNotExist(t *testing.T) {
	os.Remove("../../test")
	assert.Equal(t, true, createDirIfNotExist("../../test"), "Folder Creation should be possible")
	assert.Equal(t, false, createDirIfNotExist("?*:|"), "Directory should not be creatable")

}

func TestCreateLogfileIfNotExist(t *testing.T) {
	logFile, fileErr := createLogfileIfNotExist("../../log", "test")
	if fileErr != nil {
		t.Error("LogFile should be creatable")
	}
	if logFile == nil {
		t.Error("LogFile should be creatable")
	}
}

func TestCloseLogfile(t *testing.T) {

	file, fileErr := createLogfileIfNotExist("../../log", "test")
	if fileErr != nil {
		t.Error("LogFile is not creatable")
	}
	assert.Equal(t, true, closeLogFile(file), "LogFile Closing should be possible")

}

func TestShutdownLogging(t *testing.T) {
	assert.Equal(t, true, ShutdownLogging(), "Not all Logs could be closed gracefully")
}
