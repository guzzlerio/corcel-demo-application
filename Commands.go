package main

const (
	//NameCreateTicketCommand ...
	NameCreateTicketCommand = "CreateTicketCommand"

	//NameDeleteTicketCommand ...
	NameDeleteTicketCommand = "DeleteTicketCommand"

	//NameUpdateTicketCommand ...
	NameUpdateTicketCommand = "UpdateTicketCommand"

	//NameCloseTicketCommand ...
	NameCloseTicketCommand = "CloseTicketCommand"

	//NameAddTicketCommentCommand ...
	NameAddTicketCommentCommand = "AddTicketCommentCommand"
)

//CreateTicketCommand ...
type CreateTicketCommand struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

//UpdateTicketCommand ...
type UpdateTicketCommand struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

//CloseTicketCommand ...
type CloseTicketCommand struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
	Status  int    `json:"status"`
}

//DeleteTicketCommand ...
type DeleteTicketCommand struct {
	ID string `json:"id"`
}

//AddTicketCommentCommand ...
type AddTicketCommentCommand struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
}
