package main_test

import (
	"testing"

	. "github.com/guzzlerio/corcel-demo-application"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommandHandlers(t *testing.T) {

	var store = CreateInProcessTicketStore()

	Convey("CreateTicketCommandHandler", t, func() {
		Convey("Handles", func() {

			var handler = NewCreateTicketCommandHandler(store)

			var command = CreateTicketCommand{
				ID:    "1",
				Title: "Some Title",
				Body:  "Body",
			}

			handler.Handle(command)

			var ticket, err = store.GetByID("1")

			So(err, ShouldEqual, nil)
			So(ticket.ID, ShouldEqual, command.ID)
			So(ticket.Title, ShouldEqual, command.Title)
			So(ticket.Body, ShouldEqual, command.Body)
			So(ticket.Status, ShouldEqual, StatusTicketOpen)
		})
	})

	Convey("UpdateTicketCommandHandler", t, func() {

		Convey("Handles", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1", Status: StatusTicketOpen})
			var handler = NewUpdateTicketCommandHandler(store)

			var command = UpdateTicketCommand{
				ID:    "1",
				Title: "A Updated",
				Body:  "Body 1 Updated",
			}

			handler.Handle(command)

			var ticket, err = store.GetByID("1")
			So(err, ShouldEqual, nil)
			So(ticket.ID, ShouldEqual, command.ID)
			So(ticket.Title, ShouldEqual, command.Title)
			So(ticket.Body, ShouldEqual, command.Body)
			So(ticket.Status, ShouldEqual, StatusTicketOpen)
		})

		Convey("Returns an error", func() {

			var handler = NewUpdateTicketCommandHandler(store)

			var command = UpdateTicketCommand{
				ID:    "2",
				Title: "A Updated",
				Body:  "Body 1 Updated",
			}

			var err = handler.Handle(command)
			So(err, ShouldEqual, ErrNoTicketFound)
		})
	})

	Convey("DeleteTicketCommandHandler", t, func() {
		var handler = NewDeleteTicketCommandHandler(store)

		Convey("Handles", func() {

			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

			var command = DeleteTicketCommand{
				ID: "1",
			}

			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, nil)

			var _, err = store.GetByID("1")

			So(err, ShouldEqual, ErrNoTicketFound)

			var result = store.Query()

			So(len(result), ShouldEqual, 0)
		})

		Convey("Handles returns an error when a ticket with ID does not exist", func() {
			var command = DeleteTicketCommand{
				ID: "1",
			}

			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, ErrNoTicketFound)
		})
	})

	Convey("AddTicketCommentCommandHandler", t, func() {
		var handler = NewAddTicketCommentCommandHandler(store)

		Convey("Handles", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

			var command = AddTicketCommentCommand{
				ID:      "1",
				Comment: "Something to say",
			}

			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, nil)

			var ticket, err = store.GetByID("1")

			So(err, ShouldEqual, nil)

			So(len(ticket.Comments), ShouldEqual, 1)
		})

		Convey("Returns error", func() {
			var command = AddTicketCommentCommand{
				ID:      "2",
				Comment: "Something to say",
			}

			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, ErrNoTicketFound)
		})
	})

	Convey("CloseTicketCommandHandler", t, func() {
		var handler = NewCloseTicketCommandHandler(store)
		store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

		Convey("Sets the status", func() {
			var command = CloseTicketCommand{
				ID:     "1",
				Status: StatusTicketComplete,
			}
			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, nil)

			var result, _ = store.GetByID("1")
			So(result.Status, ShouldEqual, StatusTicketComplete)
		})

		Convey("Adds the comment to the list", func() {
			var command = CloseTicketCommand{
				ID:      "1",
				Status:  StatusTicketComplete,
				Comment: "The work has finished",
			}
			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, nil)

			var result, _ = store.GetByID("1")
			So(len(result.Comments), ShouldEqual, 1)
		})

		Convey("Returns error", func() {
			var command = CloseTicketCommand{
				ID:      "2",
				Status:  StatusTicketComplete,
				Comment: "The work has finished",
			}
			var handleError = handler.Handle(command)

			So(handleError, ShouldEqual, ErrNoTicketFound)
		})
	})
}
