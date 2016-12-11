package main_test

import (
	"testing"

	. "github.com/guzzlerio/corcel-demo-application"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInfrastructure(t *testing.T) {

	Convey("InProcessTicketStore", t, func() {
		var store = CreateInProcessTicketStore()

		Convey("Save adds a new ticket to the store", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

			var ticket, err = store.GetByID("1")

			So(err, ShouldEqual, nil)
			So(ticket.Title, ShouldEqual, "A")
		})

		Convey("Save updates an existing ticket", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})
			store.Save(Ticket{ID: "1", Title: "B", Body: "Body 1"})

			var ticket, err = store.GetByID("1")

			So(err, ShouldEqual, nil)
			So(ticket.Title, ShouldEqual, "B")
		})

		Convey("Returns all tickets which are not marked as deleted", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})
			store.Save(Ticket{ID: "2", Title: "B", Body: "Body 2"})
			store.Save(Ticket{ID: "3", Title: "C", Body: "Body 3", Status: StatusTicketDeleted})

			var result = store.Query()

			So(len(result), ShouldEqual, 2)
		})

		Convey("Return ErrNoTicketID", func() {
			var err = store.Save(Ticket{})

			So(err, ShouldEqual, ErrNoTicketID)
		})

	})
}
