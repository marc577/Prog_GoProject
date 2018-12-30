//Matrikelnummern:
//9188103
//1798794
//4717960
package storagehandler

import (
	"io/ioutil"
	"log"
	"os"
)

func writeJSONToFile(file string, jObj []byte) bool {
	os.Remove(file)
	err := ioutil.WriteFile(file, jObj, 0666)
	if err == nil {
		return true
	}
	return false
}

func readJSONFromFile(file string) []byte {

	// https://tutorialedge.net/golang/parsing-json-with-golang/

	jsonFile, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
