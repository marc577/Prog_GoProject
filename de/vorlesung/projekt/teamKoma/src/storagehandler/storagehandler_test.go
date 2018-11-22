package storagehandler

import (
	"testing"
)

var testUserStorageFile = "../../storage/users.json"
var testTicketStorageDir = "../../storage/tickets/"

/* ************************************
** TICKET TEST FUNCTIONS
************************************ */
func TestTicketHandling(t *testing.T) {

	// Load all tickets from storage
	var storageHandler = New(testUserStorageFile, testTicketStorageDir)
	var allTickets = *storageHandler.GetTickets()

	// Check if allTickets is nil
	if allTickets == nil {
		t.Error("Error in allTickets")
	}
	var originLen = len(allTickets)

	// Check if the subject is correct after creation
	var testTicket = storageHandler.CreateTicket("TestSubject", "TestMail", "TestText")
	if testTicket.Subject != "TestSubject" {
		t.Error("Ticket subject is wrong")
	}

	// Check if the scope vaiable was updated after creation
	allTickets = *storageHandler.GetTickets()
	if (originLen + 1) != len(allTickets) {
		t.Error("Ticket is not add in scope variable")
	}

	// Check if the scope variable has changed after update an item
	var openTicketsLen = len(*storageHandler.GetOpenTickets())
	testTicket.SetTicketStateClosed()
	var newOpenTicketLen = len(*storageHandler.GetOpenTickets())
	if openTicketsLen != (newOpenTicketLen + 1) {
		t.Error("Ticket is not up to date in scope variable")
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
	var storageHandler = New(testUserStorageFile, testTicketStorageDir)

	var userName = "SuperTestUser2"
	var userPassword = "SuperPasswort"
	var usersLen = len(*storageHandler.GetUsers())
	if storageHandler.CreateUser(userName, userPassword) == false {
		t.Error("user could not be created")
	}
	var newUsersLen = len(*storageHandler.GetUsers())

	if usersLen != newUsersLen-1 {
		t.Error("User is not updated in scope variable")
	}

	if storageHandler.CreateUser(userName, userPassword) {
		t.Error("User is duplicated")
	}

	if storageHandler.VerifyUser(userName, userPassword) == false {
		t.Error("User password could not verified")
	}

	if storageHandler.VerifyUser(userName, "wrongPassword") {
		t.Error("Userpassword should be wrong")
	}

	if storageHandler.DeleteUser(userName) == false {
		t.Error("User could not deleted")
	}

	if storageHandler.DeleteUser(userName) {
		t.Error("User should not be deleted")
	}

}
