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

func (handler *StorageHandler) setTickets(newTickets []Ticket) {
	handler.tickets = newTickets
}

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

// GetNotClosedTicketsByProcessor returns an array of all open or in processing tickets by a processor
func (handler *StorageHandler) GetNotClosedTicketsByProcessor(processor string) *[]Ticket {
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

func (handler *StorageHandler) updateTicketInScopeVariable(ticket Ticket) {
	var newTickets []Ticket
	for _, t := range handler.tickets {
		if t.ID == ticket.ID {
			newTickets = append(newTickets, ticket)
		} else {
			newTickets = append(newTickets, t)
		}
	}
	handler.setTickets(newTickets)
}

// UpdateTicket updates the ticket in memory and rom
// Returns the updated Ticket
func (handler *StorageHandler) UpdateTicket(ticket Ticket) Ticket {
	// Update in memory storage
	var t = ticket.writeTicketToMemory()
	// Update in scope variable
	handler.updateTicketInScopeVariable(t)
	return ticket
}

// CreateTicket creates a new ticket on persistant storage and rom
// Returns the created Ticket
func (handler *StorageHandler) CreateTicket(subject string, email string, text string) Ticket {
	var ticket = storeTicket(handler, subject, email, text)
	handler.tickets = append(handler.tickets, ticket)
	return ticket
}

/*
func (handler *StorageHandler) deleteTicket(ticket Ticket) bool {

	var ticket = storeTicket(handler, subject, email, text)
	handler.tickets = append(handler.tickets, ticket)
	return ticket
}
*/

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
