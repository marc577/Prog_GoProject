//Matrikelnummern:
//9188103
//1798794
//4717960
package storagehandler

// Email hold information about an email to send
type Email struct {
	TicketID   string
	TicketItem TicketItem
}

// GetMailsToSend returns an array of all mails which have to bee sended
func (handler *StorageHandler) GetMailsToSend() []Email {
	var mails2send []Email
	for _, ticket := range handler.tickets {
		for _, item := range ticket.Items {
			if item.IsToSend && !item.IsSended {
				mails2send = append(mails2send, Email{ticket.ID, item})
			}
		}
	}
	return mails2send
}

// SetSentMails sets the status to sent of the ticketItems
func (handler *StorageHandler) SetSentMails(sendedMails []Email) bool {
	for _, email := range sendedMails {
		var ticket, error = handler.GetTicketByID(email.TicketID)
		if error != nil {
			return false
		}
		for _, item := range ticket.Items {
			if item.CreationDate == email.TicketItem.CreationDate {
				item.IsSended = true
				ticket.Items[item.CreationDate] = item
				ticket, error = handler.UpdateTicket(ticket)
				if error != nil {
					return false
				}
			}
		}
	}
	return true
}
