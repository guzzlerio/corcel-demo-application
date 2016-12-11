package main

//InProcessTicketStore ...
type InProcessTicketStore struct {
	tickets []Ticket
}

//Query ...
func (instance *InProcessTicketStore) Query() []Ticket {
	var returnTickets = []Ticket{}
	for _, ticket := range instance.tickets {
		if ticket.Status != StatusTicketDeleted {
			returnTickets = append(returnTickets, ticket)
		}
	}

	return returnTickets
}

//CreateInProcessTicketStore ...
func CreateInProcessTicketStore() *InProcessTicketStore {
	return &InProcessTicketStore{
		tickets: []Ticket{},
	}
}

//GetByID ...
func (instance *InProcessTicketStore) GetByID(id string) (Ticket, error) {
	for _, ticket := range instance.tickets {
		if ticket.ID == id && ticket.Status != StatusTicketDeleted {
			return ticket, nil
		}
	}
	return Ticket{}, ErrNoTicketFound
}

//Save ...
func (instance *InProcessTicketStore) Save(ticketToSave Ticket) error {
	if ticketToSave.ID == "" {
		return ErrNoTicketID
	}

	for index, ticket := range instance.tickets {
		if ticket.ID == ticketToSave.ID {
			instance.tickets[index] = ticketToSave
			return nil
		}
	}
	instance.tickets = append(instance.tickets, ticketToSave)
	return nil
}
