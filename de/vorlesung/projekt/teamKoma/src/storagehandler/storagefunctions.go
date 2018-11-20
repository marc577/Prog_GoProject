package storagehandler

import (
	"fmt"
	"io/ioutil"
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
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened users.json")
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}