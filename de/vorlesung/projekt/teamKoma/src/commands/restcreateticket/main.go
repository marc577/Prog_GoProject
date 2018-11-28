package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
)

type email struct {
	Mail        string `json:"mail"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

func main() {
	host := flag.String("host", "127.0.0.1", "Ticketsystem Hostname")
	port := flag.Int("port", 8443, "Ticksetsystem Webserver Port")
	mail := flag.String("mail", "test@test.de", "Ticketsender mail adress")
	subject := flag.String("subject", "testsubject", "Ticket subject")
	description := flag.String("desc", "testdesc", "Ticket description")
	flag.Parse()
	jsonData := genJSONData(*mail, *subject, *description)
	sendReq(*host, *port, jsonData)
}

func genJSONData(mail string, subject string, desc string) []byte {
	mail2send := email{
		Mail:        mail,
		Subject:     subject,
		Description: desc}
	jsonData, err := json.Marshal(mail2send)
	if err != nil {
		log.Fatal("Error serializing Object", err)
	}
	return jsonData
}

func sendReq(host string, port int, jsonData []byte) {
	res, err := http.Post("https://"+host+":"+strconv.Itoa(port)+"/api/new", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error connecting to Server", err)
	}
	if res.StatusCode == http.StatusOK {
		log.Println("Sent Request to API")
	}
}
