package main

//CreateTicketCommandHandler ...
type CreateTicketCommandHandler struct {
	store TicketStore
}

//NewCreateTicketCommandHandler ...
func NewCreateTicketCommandHandler(store TicketStore) CreateTicketCommandHandler {
	return CreateTicketCommandHandler{
		store: store,
	}
}

//Handle ...
func (instance CreateTicketCommandHandler) Handle(command interface{}) error {
	createTicketCommand := command.(CreateTicketCommand)

	var ticket = Ticket{
		ID:       createTicketCommand.ID,
		Title:    createTicketCommand.Title,
		Body:     createTicketCommand.Body,
		Comments: []TicketComment{},
		Status:   StatusTicketOpen,
	}

	instance.store.Save(ticket)

	return nil
}

//CloseTicketCommandHandler ...
type CloseTicketCommandHandler struct {
	store TicketStore
}

//NewCloseTicketCommandHandler ...
func NewCloseTicketCommandHandler(store TicketStore) CloseTicketCommandHandler {
	return CloseTicketCommandHandler{
		store: store,
	}
}

//Handle ...
func (instance CloseTicketCommandHandler) Handle(command interface{}) error {
	closeTicketCommand := command.(CloseTicketCommand)

	ticket, err := instance.store.GetByID(closeTicketCommand.ID)

	if err != nil {
		return err
	}

	return instance.store.Save(ticket.
		Comment(closeTicketCommand.Comment).
		ChangeStatus(closeTicketCommand.Status))
}

//AddTicketCommentCommandHandler ...
type AddTicketCommentCommandHandler struct {
	store TicketStore
}

//NewAddTicketCommentCommandHandler ...
func NewAddTicketCommentCommandHandler(store TicketStore) AddTicketCommentCommandHandler {
	return AddTicketCommentCommandHandler{
		store: store,
	}
}

//Handle ...
func (instance AddTicketCommentCommandHandler) Handle(command interface{}) error {
	addTicketCommentCommand := command.(AddTicketCommentCommand)

	ticket, err := instance.store.GetByID(addTicketCommentCommand.ID)

	if err != nil {
		return err
	}

	return instance.store.Save(ticket.
		Comment(addTicketCommentCommand.Comment))
}

//UpdateTicketCommandHandler ...
type UpdateTicketCommandHandler struct {
	store TicketStore
}

//NewUpdateTicketCommandHandler ...
func NewUpdateTicketCommandHandler(store TicketStore) UpdateTicketCommandHandler {
	return UpdateTicketCommandHandler{
		store: store,
	}
}

//Handle ...
func (instance UpdateTicketCommandHandler) Handle(command interface{}) error {
	updateTicketCommand := command.(UpdateTicketCommand)

	ticket, err := instance.store.GetByID(updateTicketCommand.ID)

	if err != nil {
		return err
	}

	ticket.Title = updateTicketCommand.Title
	ticket.Body = updateTicketCommand.Body

	return instance.store.Save(ticket)
}

//DeleteTicketCommandHandler ...
type DeleteTicketCommandHandler struct {
	store TicketStore
}

//NewDeleteTicketCommandHandler ...
func NewDeleteTicketCommandHandler(store TicketStore) DeleteTicketCommandHandler {
	return DeleteTicketCommandHandler{
		store: store,
	}
}

//Handle ...
func (instance DeleteTicketCommandHandler) Handle(command interface{}) error {
	deleteTicketCommand := command.(DeleteTicketCommand)

	ticket, err := instance.store.GetByID(deleteTicketCommand.ID)
	if err != nil {
		return err
	}
	instance.store.Save(ticket.Delete())

	return nil
}
