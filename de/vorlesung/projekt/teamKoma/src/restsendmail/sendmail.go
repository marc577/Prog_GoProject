package restsendmail

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

func main() {

	host := flag.String("host", "127.0.0.1", "Ticketsystem Hostname")
	port := flag.Int("port", 8443, "Ticksetsystem Webserver Port")

	flag.Parse()
	log.Println("Flags parsed: Host:" + *host)
	log.Println("Flags Parsed: Port:" + strconv.Itoa(*port))

	fmt.Println("Client configured to connect to " + *host + ":" + strconv.Itoa(*port))
	fmt.Println("Commands to execute:")
	fmt.Println("Catch Mails to be sent: 1")
	fmt.Println("Exit : exit")

	cmd := "0"
	fmt.Sscanln(cmd)

	for cmd != "1" || cmd != "exit" {
		switch cmd {
		case "1": //do shit function grabbing mail
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Enter a correct command!")
		}
	}

}
