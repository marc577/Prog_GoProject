//Matrikelnummern:
//9188103
//1798794
//4717960
package storagehandler

import (
	"testing"
	"time"
)

var testUserStorageFile = "../../storage/users.json"
var testTicketStorageDir = "../../storage/tickets/"

var userName = "SuperTestUser2"
var userPassword = "SuperPasswort"

/* ************************************
** TICKET TEST FUNCTIONS
************************************ */

func TestTicketItems(t *testing.T) {
	var storageHandler = New(testUserStorageFile, testTicketStorageDir)
	var testTicket, _ = storageHandler.CreateTicket("TestSubject", "First Entry", "TestMail", "TestName")
	time.Sleep(1 * time.Second)
	testTicket.AddEntry2Ticket("TestCreator", "second entry", false, "TestEmailTo")
	time.Sleep(1 * time.Second)
	testTicket.AddEntry2Ticket("TestCreator", "third Entry", false, "TestMailTo")
	time.Sleep(1 * time.Second)
	testTicket.AddEntry2Ticket("TestCreator", "last entry", false, "TestEmailTo")
	time.Sleep(1 * time.Second)
	ticketEntry, error := testTicket.GetLastEntryOfTicket()
	if error != nil && ticketEntry != (TicketItem{}) {
		t.Error("Error")
	}

	ticketEntry, error = testTicket.GetFirstEntryOfTicket()
	if error != nil && ticketEntry != (TicketItem{}) {
		t.Error("Error")
	}
}
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
	var testTicket, errCreateTicket = storageHandler.CreateTicket("TestSubject", "TestText", "TestMail", "TestName")
	if testTicket.Subject != "TestSubject" && errCreateTicket != nil {
		t.Error("Ticket subject is wrong")
	}

	testTicket, errorSetSubject := testTicket.SetSubject("TestSubject2")
	if testTicket.Subject != "TestSubject2" && errorSetSubject != nil {
		t.Error("Ticket subject is wrong")
	}

	// Check if the scope vaiable was updated after creation
	allTickets = *storageHandler.GetTickets()
	if (originLen + 1) != len(allTickets) {
		t.Error("Ticket is not add in scope variable")
	}

	var newTestTicket, errGetTicketByID = storageHandler.GetTicketByID(testTicket.ID)
	if newTestTicket.ID != testTicket.ID && errGetTicketByID != nil {
		t.Error("Error by get ticket by ID")
	}

	newTestTicket, errGetTicketByID = storageHandler.GetTicketByID("")
	if errGetTicketByID == nil {
		t.Error("Error should be null")
	}

	// Check if the scope variable has changed after update an itemState to open
	var openTicketsLen = len(*storageHandler.GetOpenTickets())
	testTicket, errorSetTicketStateClosed := testTicket.SetTicketStateClosed()
	var newOpenTicketLen = len(*storageHandler.GetOpenTickets())
	if openTicketsLen != (newOpenTicketLen+1) && errorSetTicketStateClosed != nil {
		t.Error("Ticket is not up to date in rom")
	}

	testTicket, errorSetTicketStateOpen := testTicket.SetTicketStateOpen()
	if testTicket.TicketState != TSOpen && errorSetTicketStateOpen != nil {
		t.Error("Ticket state is wrong")
	}

	// Creates an test user for set in processing by user
	if storageHandler.CreateUser(userName, userPassword) == false {
		t.Error("user could not be created")
	}

	// Check if the scope variable has changed after update an itemState to inProcessing
	var ticketsByProcessorLen = len(*storageHandler.GetInProgressTicketsByProcessor(userName))
	testTicket.SetTicketStateInProgress(userName)
	var newTicketsByProcessorLen = len(*storageHandler.GetInProgressTicketsByProcessor(userName))
	if ticketsByProcessorLen != (newTicketsByProcessorLen - 1) {
		t.Error("Ticket is not up to date in rom")
	}

	if storageHandler.DeleteUser(userName) == false {
		t.Error("User could not deleted")
	}

	var ticketEntryLen = len(testTicket.Items)
	testTicket.AddEntry2Ticket("TestCreator", "An entry", false, "TestEmailTo")
	var newTicketEntryLen = len(testTicket.Items)
	if ticketEntryLen != newTicketEntryLen-1 {
		t.Error("Error while adding ticket entry to testticket")
	}

	// Test combine Tickets
	var testTicketCombine1, errc1 = storageHandler.CreateTicket("TestCombine1", "TestText", "TestMail", "TestNameC1")
	var testTicketCombine2, errc2 = storageHandler.CreateTicket("TestCombine2", "TestText", "TestMail", "TestNameC2")
	if errc1 != nil && errc2 != nil {
		t.Error("TestCombineTicketsCouldNotBeCreated")
	}
	testTicketCombine1.SetTicketStateInProgress("Processor")
	testTicketCombine2.SetTicketStateInProgress("Processor")

	storageHandler.CombineTickets(testTicketCombine1.ID, testTicketCombine2.ID)

	// Test mails to send
	var testTicketMail2send, errSent = storageHandler.CreateTicket("TestTicket2Mail", "TestText", "TestMail", "TestName")
	if errSent != nil {
		t.Error("Error in create test ticket for sending mail")
	}
	testTicketMail2send.AddEntry2Ticket("TestCreator2Send", "TestText2send", true, "testMailTo")
	var tickets2Send = storageHandler.GetMailsToSend()
	var lenTickets2Send = len(tickets2Send)
	var ticketsWhichWillSend []Email
	ticketsWhichWillSend = append(ticketsWhichWillSend, tickets2Send[0])
	storageHandler.SetSentMails(ticketsWhichWillSend)
	tickets2Send = storageHandler.GetMailsToSend()
	var newLenTickets2Send = len(tickets2Send)
	if newLenTickets2Send == lenTickets2Send {
		t.Error("Error while sending mails")
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

	var usersLen = len(*storageHandler.GetUsers())
	var availableUsersLen = len(*storageHandler.GetAvailableUsers())
	if storageHandler.CreateUser(userName, userPassword) == false {
		t.Error("user could not be created")
	}
	var newUsersLen = len(*storageHandler.GetUsers())
	var newAvailableUsersLen = len(*storageHandler.GetAvailableUsers())

	if usersLen != newUsersLen-1 {
		t.Error("User is not updated in scope variable")
	}

	if availableUsersLen != newAvailableUsersLen-1 {
		t.Error("User might not created correkt")
	}

	if storageHandler.CreateUser(userName, userPassword) {
		t.Error("User is duplicated")
	}

	var oldHolidayState = storageHandler.GetUserByUserName(userName).HasHoliday
	storageHandler.ToggleHoliday(userName)
	var newHolidayState = storageHandler.GetUserByUserName(userName).HasHoliday

	if oldHolidayState == newHolidayState {
		t.Error("User holiday not toggled correct")
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
