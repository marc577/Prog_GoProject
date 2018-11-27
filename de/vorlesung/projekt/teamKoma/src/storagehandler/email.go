package storagehandler

// Email hold information about an email to send
type Email struct {
	ticketID   string
	ticketItem TicketItem
}

// GetMailsToSend returns an array of all mails which have to bee sended
func (handler *StorageHandler) GetMailsToSend() []Email {
	var mails2send []Email
	for _, ticket := range handler.tickets {
		for _, item := range ticket.Items {
			if item.IsToSend {
				mails2send = append(mails2send, Email{ticket.ID, item})
			}
		}
	}
	return mails2send
}

// SetSendedMails sets the status to sendet of the ticketItems
func (handler *StorageHandler) SetSendedMails(sendedMails []Email) bool {
	for _, email := range sendedMails {
		var ticket, error = handler.GetTicketByID(email.ticketID)
		if error != nil {
			return false
		}
		for _, item := range ticket.Items {
			if item.CreationDate == email.ticketItem.CreationDate {
				item.IsSended = true
				ticket, error = handler.UpdateTicket(ticket)
				if error != nil {
					return false
				}
			}
		}
	}
	return true
}
