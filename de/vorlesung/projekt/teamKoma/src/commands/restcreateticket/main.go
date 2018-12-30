//Matrikelnummern:
//9188103
//1798794
//4717960
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
)

type email struct {
	Mail        string `json:"mail"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

func main() {
	//set flags
	host := flag.String("host", "127.0.0.1", "Ticketsystem Hostname")
	port := flag.Int("port", 8443, "Ticksetsystem Webserver Port")
	mail := flag.String("mail", "test@test.de", "Ticketsender mail address")
	firstname := flag.String("firstname", "test", "Ticketsender firstname")
	lastname := flag.String("lastname", "test", "Ticketsender lastname")
	subject := flag.String("subject", "testsubject", "Ticket subject")
	description := flag.String("desc", "testdesc", "Ticket description")
	flag.Parse()

	//generate JSON Data from Flags
	jsonData := genJSONData(*mail, *firstname, *lastname, *subject, *description)

	//send Request to API
	sendReq(*host, *port, jsonData)
}

//genJSONData generates valid JSON Data from given input
func genJSONData(mail string, firstname string, lastname string, subject string, desc string) []byte {
	mail2send := email{
		Mail:        mail,
		Firstname:   firstname,
		Lastname:    lastname,
		Subject:     subject,
		Description: desc}
	jsonData, err := json.Marshal(mail2send)
	if err != nil {
		log.Fatal("Error serializing Object", err)
	}
	return jsonData
}

//sendReq sends the POST Request to the API
func sendReq(host string, port int, jsonData []byte) error {
	//need to ignore unknown cert authority for Request
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	res, err := http.Post("https://"+host+":"+strconv.Itoa(port)+"/api/new", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error connecting to Server", err)
		return err
	}
	if res.StatusCode == http.StatusOK {
		log.Println("Sent Request to API")
	}
	return nil
}
