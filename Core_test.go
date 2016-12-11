package main_test

import (
	"errors"
	"testing"

	. "github.com/guzzlerio/corcel-demo-application"
	. "github.com/smartystreets/goconvey/convey"
)

type TestCommand struct{}

type TestHandler struct {
	err    error
	called bool
}

func (instance *TestHandler) SetError(err error) {
	instance.err = err
}

func (instance *TestHandler) Handle(command interface{}) error {
	instance.called = true
	return instance.err
}

func (instance *TestHandler) WasCalled() bool {
	return instance.called
}

func TestCore(t *testing.T) {

	Convey("Ticket", t, func() {
		var ticket = Ticket{}
		Convey("Delete", func() {
			ticket = ticket.Delete()
			So(ticket.Status, ShouldEqual, StatusTicketDeleted)
		})

	})

	Convey("CommandController", t, func() {
		var controller = CommandController{
			HandlerMappings: map[string]CommandHandler{},
		}
		var testHandler = &TestHandler{}
		var testCommandName = "test"

		controller.AddHandler(testCommandName, testHandler)

		Convey("Returns ErrNoCommandHandler", func() {
			var result = controller.Handle("doesNotExist", TestCommand{})

			So(result, ShouldEqual, ErrNoCommandHandler)
		})

		Convey("Invokes the mapped handler", func() {
			var result = controller.Handle(testCommandName, TestCommand{})

			So(result, ShouldEqual, nil)
			So(testHandler.WasCalled(), ShouldEqual, true)
		})

		Convey("Returns the error from a command handler", func() {
			var expectedError = errors.New("")
			testHandler.SetError(expectedError)
			var result = controller.Handle(testCommandName, TestCommand{})

			So(result, ShouldEqual, expectedError)
			So(testHandler.WasCalled(), ShouldEqual, true)
		})
	})
}
