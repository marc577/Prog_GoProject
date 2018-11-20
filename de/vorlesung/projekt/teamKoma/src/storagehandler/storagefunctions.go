package storagehandler

import (
	"io/ioutil"
	"logging"
	"os"
)

func writeJSONToFile(file string, jObj []byte) bool {
	err := ioutil.WriteFile(file, jObj, 0777)
	if err == nil {
		return true
	}
	return false
}

func readJSONFromFile(file string) []byte {

	// https://tutorialedge.net/golang/parsing-json-with-golang/

	jsonFile, err := os.Open(file)

	if err != nil {
		logging.Error.Panic(err)
	} else {
		logging.Info.Println("Successfully Opened users.json")
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
