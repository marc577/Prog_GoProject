package restsendmail

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TicketState int

const (
	// TSOpen represents the open state
	TSOpen TicketState = 0
	// TSInProgress represents the in process state
	TSInProgress TicketState = 1
	// TSClosed represents the closed state
	TSClosed TicketState = 2
)

// TicketItem represents an entry of a ticket
type TicketItem struct {
	CreationDate time.Time `json:"creationDate"`
	Email        string    `json:"email"`
	Text         string    `json:"text"`
}

// Ticket represents a ticket
type Ticket struct {
	ID          string                   `json:"id"`
	Subject     string                   `json:"subject"`
	TicketState TicketState              `json:"ticketState"`
	Processor   string                   `json:"processor"`
	Items       map[time.Time]TicketItem `json:"items"`
}

func main() {

	host := flag.String("host", "127.0.0.1", "Ticketsystem Hostname")
	port := flag.Int("port", 8443, "Ticksetsystem Webserver Port")
	user := flag.String("user", "dummy", "Your Ticketsystem Username")
	pass := flag.String("password", "dummy", "Your Ticketsystem Password")

	flag.Parse()
	log.Println("Flags parsed: Host:" + *host)
	log.Println("Flags Parsed: Port:" + strconv.Itoa(*port))
	log.Println("Flags parsed: User:" + *user)
	log.Println("Flags parsed: Password:" + *pass)

	fmt.Println("Client configured to connect to " + *host + ":" + strconv.Itoa(*port))
	fmt.Println("Commands to execute:")
	fmt.Println("Catch Mails to be sent: 1")
	fmt.Println("Exit : exit")

	cmd := "0"
	fmt.Sscanln(cmd)

	for cmd != "1" || cmd != "exit" {
		switch cmd {
		case "1":
			var []Ticket:=grabMailsToSend(*host, *port, *user, *pass)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Enter a correct command!")
		}
	}

}

func grabMailsToSend(host string, port int, user string, pass string) []Ticket {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://"+host+":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, pass)
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	tickets := make([]Ticket, 0)
	json.Unmarshal(body, &tickets)
	return tickets
}

func setSentFlag() {

}
