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
	ID          string                   `json:"id"`
	Subject     string                   `json:"subject"`
	TicketState TicketState              `json:"ticketState"`
	Processor   string                   `json:"processor"`
	Items       map[time.Time]TicketItem `json:"items"`
}

// SetSubject sets the subject in the given ticket
func (ticket Ticket) SetSubject(subject string) Ticket {
	ticket.Subject = subject
	updateTicketInScopeVariable(ticket)
	return ticket.writeTicketToMemory()
}

// SetTicketStateOpen sets the subject in the given ticket
// The processor will be resetet
func (ticket Ticket) SetTicketStateOpen() Ticket {
	ticket.TicketState = TSOpen
	ticket.Processor = ""
	return UpdateTicket(ticket)
}

// SetTicketStateInProgress sets the subject in the given ticket
func (ticket Ticket) SetTicketStateInProgress(processor string) Ticket {
	ticket.TicketState = TSInProgress
	ticket.Processor = processor
	return UpdateTicket(ticket)
}

// SetTicketStateClosed sets the subject in the given ticket
// The processor will be resetet
func (ticket Ticket) SetTicketStateClosed() Ticket {
	ticket.TicketState = TSClosed
	ticket.Processor = ""
	return UpdateTicket(ticket)
}

// AddEntry2Ticket adds an entry to the given ticket
func (ticket Ticket) AddEntry2Ticket(email string, text string) Ticket {
	currTime := time.Now()
	ticket.Items[currTime] = TicketItem{currTime, email, text}
	return UpdateTicket(ticket)
}

// Delete the ticket
func (ticket Ticket) Delete() {
	// TODO: other solution
	var newTickets []Ticket
	var oldTickets = GetTickets()
	for i := 0; i < len(*oldTickets); i++ {
		if (*oldTickets)[i].ID != ticket.ID {
			newTickets = append(newTickets, (*oldTickets)[i])
		}
	}
	setTickets(newTickets)
}

func loadSpecificTicketFromMemory(ticketID string) Ticket {
	var ticket Ticket
	var byteValue = readJSONFromFile(ticketStoreDir + ticketID + ".json")
	json.Unmarshal(byteValue, &ticket)
	return ticket
}

func loadFilesFromMemory() []Ticket {
	var tickets []Ticket

	file, err := os.Open(ticketStoreDir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		if strings.Contains(name, ".json") {
			var ticket Ticket
			var byteValue = readJSONFromFile(ticketStoreDir + name)
			json.Unmarshal(byteValue, &ticket)
			tickets = append(tickets, ticket)
		}
	}
	return tickets
}

func (ticket Ticket) writeTicketToMemory() Ticket {
	result, err := json.Marshal(ticket)
	if err != nil {
		fmt.Println("Error while add user")
	}
	if writeJSONToFile((ticketStoreDir+ticket.ID+".json"), result) == false {
		fmt.Println("Error while write Ticket to memory")
	}
	return ticket
}

// storeTicket writes a new json-File to the memory with in a ticket
func storeTicket(subject string, email string, text string) Ticket {
	currentTime := time.Now()
	ticketID := string(currentTime.Format("20060102150405")) + "_" + email
	item := TicketItem{currentTime, email, text}
	mItems := make(map[time.Time]TicketItem)
	mItems[currentTime] = item
	newTicket := &Ticket{ticketID, subject, TSOpen, "", mItems}
	return newTicket.writeTicketToMemory()
}
