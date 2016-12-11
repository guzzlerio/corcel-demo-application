package main

import "github.com/uber-go/zap"

//LoggingTicketStore ...
type LoggingTicketStore struct {
	WrappedTicketStore TicketStore
}

//Query ...
func (instance *LoggingTicketStore) Query() []Ticket {
	var result = instance.WrappedTicketStore.Query()
	logger.Info("Query")
	return result
}

//GetByID ...
func (instance *LoggingTicketStore) GetByID(id string) (Ticket, error) {
	var result, err = instance.WrappedTicketStore.GetByID(id)

	if err != nil {
		logger.Error("GetTicketByID",
			zap.String("id", id),
			zap.String("error", err.Error()),
		)
		return result, err
	}

	logger.Info("GetTicketByID",
		zap.String("id", id),
	)
	return result, err
}

//Save ...
func (instance *LoggingTicketStore) Save(ticketToSave Ticket) error {
	var err = instance.WrappedTicketStore.Save(ticketToSave)
	if err != nil {
		logger.Error("Save",
			zap.String("id", ticketToSave.ID),
			zap.String("error", err.Error()),
		)
		return err
	}

	logger.Info("Save",
		zap.String("id", ticketToSave.ID),
		zap.Int("status", ticketToSave.Status),
	)
	return err
}
