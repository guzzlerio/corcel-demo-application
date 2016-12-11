package main_test

import (
	"bytes"
	"testing"

	. "github.com/guzzlerio/corcel-demo-application"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/uber-go/zap"
)

func TestLogging(t *testing.T) {

	defer InitializeLogger()

	Convey("LoggingCommandController", t, func() {
		buf := new(bytes.Buffer)

		SetLogger(zap.New(
			zap.NewJSONEncoder(), // drop timestamps in tests
			zap.DebugLevel,
			zap.Output(zap.AddSync(buf)),
		))

		var store = LoggingTicketStore{
			WrappedTicketStore: CreateInProcessTicketStore(),
		}

		Convey("Save adds a new ticket to the store", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

			So(buf.String(), ShouldContainSubstring, "\"Save\"")
		})

		Convey("Get a ticket by ID", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

			store.GetByID("1")

			So(buf.String(), ShouldContainSubstring, "\"GetTicketByID\"")
		})

		Convey("Get a ticket by ID logs error", func() {
			store.GetByID("1")

			So(buf.String(), ShouldContainSubstring, "\"GetTicketByID\"")
			So(buf.String(), ShouldContainSubstring, "\"error\"")
		})

		Convey("Save logs an error", func() {
			store.Save(Ticket{})

			So(buf.String(), ShouldContainSubstring, "\"Save\"")
			So(buf.String(), ShouldContainSubstring, "\"error\"")
		})

		Convey("Returns all tickets which are not marked as deleted", func() {
			store.Save(Ticket{ID: "1", Title: "A", Body: "Body 1"})

			store.Query()

			So(buf.String(), ShouldContainSubstring, "\"Query\"")
		})

	})
}
