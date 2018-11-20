package storagehandler

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init()
}

/* ************************************
** TICKET TEST FUNCTIONS
************************************ */
func TestTicketHandling(t *testing.T) {

	// Load all tickets from storage
	Init()
	var allTickets = GetTickets()

	// Check if allTickets is nil
	if allTickets == nil {
		t.Error("Error in allTickets")
	}

	// Check if the subject is correct
	var testTicket = CreateTicket("TestSubject", "TestMail", "TestText")
	if testTicket.Subject != "TestSubject" {
		t.Error("Ticket subject is wrong")
	}

	var originLen = len(allTickets)
	allTickets = GetTickets()

	// Check if the scope vaiable was updated
	if (originLen + 1) != len(allTickets) {
		t.Error("Ticket is not add in scope variable")
	}

	var openTickets = GetOpenTickets()
	var openTicketsLen = len(openTickets)

	testTicket.SetTicketStateClosed()
	var newOpenTicketLen = len(GetOpenTickets())

	if openTicketsLen != (newOpenTicketLen - 1) {
		t.Error("Ticket is not up to date in scope variable")
	}

}
func TestGetNotClosedTicketsByProcessor(t *testing.T) {
	Init()
	if GetNotClosedTicketsByProcessor("Klaus") == nil {
		t.Error("Error in function GetNotClosedTicketsByProcessor")
	}
}

func TestGetOpenTickets(t *testing.T) {
	Init()
	if GetOpenTickets() == nil {
		t.Error("Error in function GetOpenTickets")
	}
}

func TestDeleteTicket(t *testing.T) {
	Init()
	if 1 == 2 {
		t.Error("not implemented")
	}
}

/* ************************************
** Test User Functions
** Creates an user
** Check if the same user can not created
** verify user by passwort
** verify wrong passwort
** delete the created user
** check if a non existing canot deleted
************************************ */

func TestUserFunctions(t *testing.T) {

	var userName = "SuperTestUser"
	var userPassword = "SuperPasswort"
	if CreateUser(userName, userPassword) == false {
		t.Error("user could not be created")
	}

	if CreateUser(userName, userPassword) {
		t.Error("User is duplicated")
	}

	if VerifyUser(userName, userPassword) == false {
		t.Error("User password could not verified")
	}

	if VerifyUser(userName, "wrongPassword") {
		t.Error("Userpassword should be wrong")
	}

	if DeleteUser(userName) == false {
		t.Error("User could not deleted")
	}

	if DeleteUser(userName) {
		t.Error("User should not be deleted")
	}

}
