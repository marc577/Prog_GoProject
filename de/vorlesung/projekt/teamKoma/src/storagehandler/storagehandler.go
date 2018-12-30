//Matrikelnummern:
//9188103
//1798794
//4717960
package storagehandler

import (
	"errors"
)

// StorageHandler defines a struct of the storageHandler
type StorageHandler struct {
	tickets        []Ticket
	users          []User
	userStoreFile  string
	ticketStoreDir string
}

// The StorageWrapper Interface for the ticket application
type StorageWrapper interface {
	GetAvailableUsers() *[]User
	GetUserByUserName(userName string) User
	GetMailsToSend() []Email
	SetSendedMails(sendedMails []Email) bool
	GetTicketByID(id string) (Ticket, error)
	CreateTicket(subject string, text string, email string, name string) (Ticket, error)
	GetOpenTickets() *[]Ticket
	GetTickets() *[]Ticket
	VerifyUser(userName string, userPassword string) bool
	GetInProgressTicketsByProcessor(processor string) *[]Ticket
}

// New loads the storage into ROM and return a new StorageHandler Object
func New(argUserStoreFile string, argTicketStoreDir string) *StorageHandler {
	var handler StorageHandler
	handler.userStoreFile = argUserStoreFile
	handler.ticketStoreDir = argTicketStoreDir
	handler.loadTicketFilesFromMemory()
	handler.loadUserFromMemory()

	return &handler
}

/* ************************************
** TICKET FUNCTIONS
************************************ */

// GetTickets returns all tickets
func (handler *StorageHandler) GetTickets() *[]Ticket {
	return &handler.tickets
}

// GetTicketByID Returns a ticket by the given id
func (handler *StorageHandler) GetTicketByID(id string) (Ticket, error) {
	for _, ticket := range handler.tickets {
		if ticket.ID == id {
			return ticket, nil
		}
	}
	return Ticket{}, errors.New("can not find ticket by the given id")
}

// GetInProgressTicketsByProcessor returns an array of all open or in processing tickets by a processor
func (handler *StorageHandler) GetInProgressTicketsByProcessor(processor string) *[]Ticket {
	var openTicketsByProcessor []Ticket
	for _, ticket := range handler.tickets {
		if ticket.TicketState != 2 && ticket.Processor == processor {
			openTicketsByProcessor = append(openTicketsByProcessor, ticket)
		}
	}
	return &openTicketsByProcessor
}

// GetOpenTickets return an array of tickets with the ticket state open
func (handler *StorageHandler) GetOpenTickets() *[]Ticket {
	var openTickets []Ticket
	for _, ticket := range handler.tickets {
		if ticket.TicketState == 0 {
			openTickets = append(openTickets, ticket)
		}
	}
	return &openTickets
}

// UpdateTicket updates the ticket in memory and rom
// Returns the updated Ticket
func (handler *StorageHandler) UpdateTicket(ticket Ticket) (Ticket, error) {
	// Update in memory storage
	ticket, error := ticket.writeTicketToMemory()
	// Update in scope variable
	for i := 0; i < len(handler.tickets); i++ {
		if handler.tickets[i].ID == ticket.ID {
			handler.tickets[i] = ticket
			break
		}
	}
	return ticket, error
}

// CreateTicket creates a new ticket on persistant storage and rom
// Returns the created Ticket
func (handler *StorageHandler) CreateTicket(subject string, text string, email string, name string) (Ticket, error) {
	var ticket, error = storeTicket(handler, subject, text, email, name)
	handler.tickets = append(handler.tickets, ticket)
	return ticket, error
}

/*
// AddEntry2Ticket adds an entry to the given ticket
func (ticket Ticket) AddEntry2Ticket(creator string, text string, isToSend bool, emailTo string) (Ticket, error) {
	currTime := time.Now()
	ticket.Items[currTime] = TicketItem{currTime, creator, text, isToSend, false, emailTo}
	return ticket.storageHandler.UpdateTicket(ticket)
}
*/

//CombineTickets combines the second ticket into the first one by given ids
func (handler *StorageHandler) CombineTickets(id1 string, id2 string) (Ticket, error) {
	var ticket1 = Ticket{}
	ticket1.ID = "-1"
	var ticket2 = Ticket{}
	ticket2.ID = "-2"
	for i := 0; i < len(handler.tickets); i++ {
		if handler.tickets[i].ID == id1 {
			ticket1 = handler.tickets[i]
		} else if handler.tickets[i].ID == id2 {
			ticket2 = handler.tickets[i]
		}
	}
	if ticket1.ID == "-1" || ticket2.ID == "-2" {
		return Ticket{}, errors.New("One of the tickets can't be found")
	}
	if ticket1.Processor != ticket2.Processor {
		return Ticket{}, errors.New("The tickets have not the same processor")
	}
	for _, item2 := range ticket2.Items {
		ticket1.Items[item2.CreationDate] = item2
	}
	var retTicket, error = handler.UpdateTicket(ticket1)
	ticket2.storageHandler.deleteTicket(ticket2)
	return retTicket, error
}

/* ************************************
** USER FUNCTIONS
************************************ */

// GetUsers returns all users
func (handler *StorageHandler) GetUsers() *[]User {
	return &handler.users
}

// GetAvailableUsers returns all users which has no holidays
func (handler *StorageHandler) GetAvailableUsers() *[]User {
	var availableUsers []User
	for _, user := range handler.users {
		if user.HasHoliday == false {
			availableUsers = append(availableUsers, user)
		}
	}
	return &availableUsers
}

// DeleteUser delets an user from memory storage
func (handler *StorageHandler) DeleteUser(userName string) bool {
	return handler.deleteUser(userName)
}

// VerifyUser check if username and password match
func (handler *StorageHandler) VerifyUser(userName string, userPassword string) bool {
	return handler.verifyUser(userName, userPassword)
}

// CreateUser create a new User
func (handler *StorageHandler) CreateUser(userName string, userPassword string) bool {
	return handler.addUser(userName, userPassword)
}
