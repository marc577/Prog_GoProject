package storagehandler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// The TicketState represents the current state of a ticket as an integer
// tsOpen = 0; tsInProgress = 1; tsClosed = 2
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
	storageHandler *StorageHandler
	ID             string                   `json:"id"`
	Subject        string                   `json:"subject"`
	TicketState    TicketState              `json:"ticketState"`
	Processor      string                   `json:"processor"`
	Items          map[time.Time]TicketItem `json:"items"`
	Email          string                   `json:"email"`
	FirstName      string                   `json:"firstName"`
	LastName       string                   `json:"lastName"`
}

// createTicketID create an id by hashing the given values
func createTicketID(currentTime time.Time, email string, FirstName string, LastName string) string {
	var id2hash = string(currentTime.Format("20060102150405")) + email + FirstName + LastName
	h := sha1.New()
	h.Write([]byte(id2hash))
	return hex.EncodeToString(h.Sum(nil))
}

// SetSubject sets the subject in the given ticket
func (ticket Ticket) SetSubject(subject string) (Ticket, error) {
	ticket.Subject = subject
	return ticket.storageHandler.UpdateTicket(ticket)
}

// SetTicketStateOpen sets the subject in the given ticket
// The processor will be resetet
func (ticket Ticket) SetTicketStateOpen() (Ticket, error) {
	ticket.TicketState = TSOpen
	ticket.Processor = ""
	return ticket.storageHandler.UpdateTicket(ticket)
}

// SetTicketStateInProgress sets the subject in the given ticket
func (ticket Ticket) SetTicketStateInProgress(processor string) (Ticket, error) {
	ticket.TicketState = TSInProgress
	ticket.Processor = processor
	return ticket.storageHandler.UpdateTicket(ticket)
}

// SetTicketStateClosed sets the subject in the given ticket
// The processor will be resetet
func (ticket Ticket) SetTicketStateClosed() (Ticket, error) {
	ticket.TicketState = TSClosed
	ticket.Processor = ""
	return ticket.storageHandler.UpdateTicket(ticket)
}

// AddEntry2Ticket adds an entry to the given ticket
func (ticket Ticket) AddEntry2Ticket(email string, text string) (Ticket, error) {
	currTime := time.Now()
	ticket.Items[currTime] = TicketItem{currTime, email, text}
	return ticket.storageHandler.UpdateTicket(ticket)
}

func (ticket Ticket) GetLastEntryOfTicket() (TicketItem, error) {
	var time time.Time
	var lastItem TicketItem
	for _, item := range ticket.Items {
		if item.CreationDate.After(time) {
			time = item.CreationDate
			lastItem = item
		}
	}
	return lastItem, nil
}

func (ticket Ticket) GetFirstEntryOfTicket() (TicketItem, error) {
	var time = time.Now()
	var firstItem TicketItem
	for _, item := range ticket.Items {
		if item.CreationDate.Before(time) {
			time = item.CreationDate
			firstItem = item
		}
	}
	return firstItem, nil
}

func (handler *StorageHandler) loadTicketFilesFromMemory() []Ticket {

	file, err := os.Open(handler.ticketStoreDir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		if strings.Contains(name, ".json") {
			var ticket Ticket
			var byteValue = readJSONFromFile(handler.ticketStoreDir + name)
			json.Unmarshal(byteValue, &ticket)
			ticket.storageHandler = handler
			handler.tickets = append(handler.tickets, ticket)
		}
	}
	return handler.tickets
}

func (ticket Ticket) writeTicketToMemory() (Ticket, error) {
	result, err := json.Marshal(ticket)
	if err != nil {
		fmt.Println("Error while add user")
	}
	if writeJSONToFile((ticket.storageHandler.ticketStoreDir+ticket.ID+".json"), result) == false {
		fmt.Println("Error while write Ticket to memory")
		return Ticket{}, errors.New("Could not write ticket to memory")
	}
	return ticket, nil
}

// storeTicket writes a new json-File to the memory with in a ticket
func storeTicket(storageHandler *StorageHandler, subject string, text string, email string, firstName string, lastName string) (Ticket, error) {
	currentTime := time.Now()
	//ticketID := string(currentTime.Format("20060102150405")) + "_" + email
	ticketID := createTicketID(currentTime, email, firstName, lastName)
	item := TicketItem{currentTime, email, text}
	mItems := make(map[time.Time]TicketItem)
	mItems[currentTime] = item
	newTicket := Ticket{storageHandler, ticketID, subject, TSOpen, "", mItems, email, firstName, lastName}
	return newTicket.writeTicketToMemory()
}
