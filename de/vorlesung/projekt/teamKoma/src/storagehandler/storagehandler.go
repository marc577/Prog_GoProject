package storagehandler

var tickets []Ticket

// Init loads the storage into ROM
func Init() bool {
	tickets = loadFilesFromMemory()
	return true
}

/* ************************************
** TICKET FUNCTIONS
************************************ */

func setTickets(newTickets []Ticket) {
	tickets = newTickets
}

// GetTickets returns all tickets
func GetTickets() []Ticket {
	return tickets
}

// GetNotClosedTicketsByProcessor returns an array of all open or in processing tickets by a processor
func GetNotClosedTicketsByProcessor(processor string) []Ticket {
	var openTicketsByProcessor []Ticket
	for _, ticket := range tickets {
		if ticket.TicketState != 2 && ticket.Processor == processor {
			openTicketsByProcessor = append(openTicketsByProcessor, ticket)
		}
	}
	return openTicketsByProcessor
}

// GetOpenTickets return an array of tickets with the ticket state open
func GetOpenTickets() []Ticket {
	var openTickets []Ticket
	for _, ticket := range tickets {
		if ticket.TicketState == 0 {
			openTickets = append(openTickets, ticket)
		}
	}
	return openTickets
}

func updateTicketInScopeVariable(ticket Ticket) {
	var newTickets []Ticket
	for _, t := range tickets {
		if t.ID == ticket.ID {
			newTickets = append(newTickets, ticket)
		} else {
			newTickets = append(newTickets, t)
		}
	}
	setTickets(newTickets)
}

// UpdateTicket updates the ticket in memory and rom
// Returns the updated Ticket
func UpdateTicket(ticket Ticket) Ticket {
	// Update in memory storage
	var t = ticket.writeTicketToMemory()
	// Update in scope variable
	updateTicketInScopeVariable(t)
	return ticket
}

// CreateTicket creates a new ticket on persistant storage and rom
// Returns the created Ticket
func CreateTicket(subject string, email string, text string) Ticket {
	var ticket = storeTicket(subject, email, text)
	tickets = append(tickets, ticket)
	return ticket
}

/* ************************************
** USER FUNCTIONS
************************************ */

// DeleteUser delets an user from memory storage
func DeleteUser(userName string) bool {
	return deleteUser(userName)
}

// VerifyUser check if username and password match
func VerifyUser(userName string, userPassword string) bool {
	return verifyUser(userName, userPassword)
}

// CreateUser create a new User
func CreateUser(userName string, userPassword string) bool {
	return addUser(userName, userPassword)
}
