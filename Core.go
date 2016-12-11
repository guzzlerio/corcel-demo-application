package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/uber-go/zap"
)

const (
	//StatusTicketOpen ...
	StatusTicketOpen = 1

	//StatusTicketDeleted ...
	StatusTicketDeleted = 2

	//StatusTicketComplete ...
	StatusTicketComplete = 3
)

var (
	//ErrNoTicketFound ...
	ErrNoTicketFound = errors.New("ErrNoTicketFound")

	//ErrNoCommandHandler ...
	ErrNoCommandHandler = errors.New("ErrNoCommandHandler")

	//ErrNoTicketID ...
	ErrNoTicketID = errors.New("ErrNoTicketID")
)

//TicketStore ...
type TicketStore interface {
	Query() []Ticket
	GetByID(id string) (Ticket, error)
	Save(ticket Ticket) error
}

//TicketComment ...
type TicketComment struct {
	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment"`
}

//Ticket ...
type Ticket struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Body     string          `json:"body"`
	Comments []TicketComment `json:"comments"`
	Status   int             `json:"status"`
}

//Delete ...
func (instance Ticket) Delete() Ticket {
	instance.Status = StatusTicketDeleted
	return instance
}

//Comment ...
func (instance Ticket) Comment(comment string) Ticket {
	var newComment = TicketComment{
		Timestamp: time.Now(),
		Comment:   comment,
	}
	instance.Comments = append(instance.Comments, newComment)
	return instance
}

//ChangeStatus ...
func (instance Ticket) ChangeStatus(status int) Ticket {
	instance.Status = status
	return instance
}

//CommandHandler ...
type CommandHandler interface {
	Handle(command interface{}) error
}

//CommandController ...
type CommandController struct {
	HandlerMappings map[string]CommandHandler
}

//AddHandler ...
func (instance CommandController) AddHandler(name string, handler CommandHandler) {
	logger.Info("AddHandler",
		zap.String("name", name),
		zap.String("type", fmt.Sprintf("%T", handler)),
	)
	instance.HandlerMappings[name] = handler
}

//Handle ...
func (instance CommandController) Handle(name string, command interface{}) error {
	if handler, ok := instance.HandlerMappings[name]; ok {
		logger.Info("Handle",
			zap.String("type", fmt.Sprintf("%T", command)),
		)
		return handler.Handle(command)
	}
	var err = ErrNoCommandHandler
	logger.Error("Handle",
		zap.String("error", err.Error()),
		zap.String("command", name),
	)
	return err
}
