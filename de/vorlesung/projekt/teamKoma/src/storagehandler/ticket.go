//Matrikelnummern:
//9188103
//1798794
//4717960
package storagehandler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	// TicketItem infos
	CreationDateString string    `json:"creationDateString"`
	CreationDate       time.Time `json:"creationDate"`
	Creator            string    `json:"creator"`
	Text               string    `json:"text"`
	// Mail infos
	IsToSend bool   `json:"isToSend"`
	IsSended bool   `json:"isSended"`
	EmailTo  string `json:"emailTo"`
}

// Ticket represents a ticket
type Ticket struct {
	storageHandler *StorageHandler
	ID             string                `json:"id"`
	Subject        string                `json:"subject"`
	TicketState    TicketState           `json:"ticketState"`
	Processor      string                `json:"processor"`
	Items          map[string]TicketItem `json:"items"`
	// Name of ticket creator
	Name string `json:"name"`
}

// createTicketID create an id by hashing the given values
func createTicketID(currentTime time.Time, email string, name string) string {
	var id2hash = string(currentTime.Format("20060102150405")) + email + name
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
	return ticket.storageHandler.UpdateTicket(ticket)
}

// AddEntry2Ticket adds an entry to the given ticket
func (ticket Ticket) AddEntry2Ticket(creator string, text string, isToSend bool, emailTo string) (Ticket, error) {
	currTime := time.Now()
	strDate := string(currTime.Format("2006010215040501"))
	ticket.Items[strDate] = TicketItem{strDate, currTime, creator, text, isToSend, false, emailTo}
	return ticket.storageHandler.UpdateTicket(ticket)
}

// GetLastEntryOfTicket returns the last ticket entry of a ticket
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

// GetFirstEntryOfTicket returns the first ticket entry of a ticket
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

// Deletes the ticket from cash and hdd
func (handler *StorageHandler) deleteTicket(ticket Ticket) {

	// delete from cache
	var i int
	for i = 0; i < len(*handler.GetTickets()); i++ {
		if (*handler.GetTickets())[i].ID == ticket.ID {
			break
		}
	}
	handler.tickets[i] = handler.tickets[len(handler.tickets)-1]
	handler.tickets[len(handler.tickets)-1] = Ticket{}
	handler.tickets = handler.tickets[:len(handler.tickets)-1]

	// Delete from hdd
	os.Remove(ticket.storageHandler.ticketStoreDir + ticket.ID + ".json")
}

func (handler *StorageHandler) loadTicketFilesFromMemory() ([]Ticket, error) {

	file, err := os.Open(handler.ticketStoreDir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	list, err := file.Readdirnames(0) // 0 to read all files and folders
	file.Close()
	if err != nil {
		return nil, err
	}
	for _, name := range list {
		if strings.Contains(name, ".json") {
			var ticket Ticket
			var byteValue = readJSONFromFile(handler.ticketStoreDir + name)
			json.Unmarshal(byteValue, &ticket)
			ticket.storageHandler = handler
			handler.tickets = append(handler.tickets, ticket)
		}
	}
	return handler.tickets, nil
}

func (ticket Ticket) writeTicketToMemory() (Ticket, error) {
	result, err := json.Marshal(ticket)
	if err != nil {
		return Ticket{}, errors.New("Could not marshal ticket")
	}
	if writeJSONToFile((ticket.storageHandler.ticketStoreDir+ticket.ID+".json"), result) == false {
		return Ticket{}, errors.New("Could not write ticket to memory")
	}
	return ticket, nil
}

// storeTicket writes a new json-File to the memory with in a ticket
func storeTicket(storageHandler *StorageHandler, subject string, text string, email string, name string) (Ticket, error) {
	currentTime := time.Now()
	strDate := string(currentTime.Format("20060102150405"))
	ticketID := createTicketID(currentTime, email, name)
	item := TicketItem{strDate, currentTime, name, text, false, false, email}
	mItems := make(map[string]TicketItem)
	mItems[strDate] = item
	newTicket := Ticket{storageHandler, ticketID, subject, TSOpen, "", mItems, name}
	return newTicket.writeTicketToMemory()
}
