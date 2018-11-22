package storagehandler

import (
	"encoding/json"
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
}

// SetSubject sets the subject in the given ticket
func (ticket Ticket) SetSubject(subject string) Ticket {
	ticket.Subject = subject
	return ticket.storageHandler.UpdateTicket(ticket)
}

// SetTicketStateOpen sets the subject in the given ticket
// The processor will be resetet
func (ticket Ticket) SetTicketStateOpen() Ticket {
	ticket.TicketState = TSOpen
	ticket.Processor = ""
	return ticket.storageHandler.UpdateTicket(ticket)
}

// SetTicketStateInProgress sets the subject in the given ticket
func (ticket Ticket) SetTicketStateInProgress(processor string) Ticket {
	ticket.TicketState = TSInProgress
	ticket.Processor = processor
	return ticket.storageHandler.UpdateTicket(ticket)
}

// SetTicketStateClosed sets the subject in the given ticket
// The processor will be resetet
func (ticket Ticket) SetTicketStateClosed() Ticket {
	ticket.TicketState = TSClosed
	ticket.Processor = ""
	return ticket.storageHandler.UpdateTicket(ticket)
}

// AddEntry2Ticket adds an entry to the given ticket
func (ticket Ticket) AddEntry2Ticket(email string, text string) Ticket {
	currTime := time.Now()
	ticket.Items[currTime] = TicketItem{currTime, email, text}
	return ticket.storageHandler.UpdateTicket(ticket)
}

// Delete the ticket
func (handler *StorageHandler) deleteTicket(argTicket Ticket) bool {

	var i int
	for i = 0; i < len(*handler.GetTickets()); i++ {
		if (*handler.GetTickets())[i].ID == argTicket.ID {
			break
		}
	}
	handler.tickets[i] = handler.tickets[len(handler.tickets)-1]
	handler.tickets[len(handler.tickets)-1] = Ticket{}
	handler.tickets = handler.tickets[:len(handler.tickets)-1]
	return true
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

func (ticket Ticket) writeTicketToMemory() Ticket {
	result, err := json.Marshal(ticket)
	if err != nil {
		fmt.Println("Error while add user")
	}
	if writeJSONToFile((ticket.storageHandler.ticketStoreDir+ticket.ID+".json"), result) == false {
		fmt.Println("Error while write Ticket to memory")
	}
	return ticket
}

// storeTicket writes a new json-File to the memory with in a ticket
func storeTicket(storageHandler *StorageHandler, subject string, email string, text string) Ticket {
	currentTime := time.Now()
	ticketID := string(currentTime.Format("20060102150405")) + "_" + email
	item := TicketItem{currentTime, email, text}
	mItems := make(map[time.Time]TicketItem)
	mItems[currentTime] = item
	newTicket := Ticket{storageHandler, ticketID, subject, TSOpen, "", mItems}
	return newTicket.writeTicketToMemory()
}
