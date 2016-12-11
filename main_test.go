package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	. "github.com/guzzlerio/corcel-demo-application"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMainAcceptance(t *testing.T) {

	var host = "localhost"
	var port = 56000
	var client = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	Convey("Acceptance Tests", t, func() {

		var server = CreateTicketServer(port, host)
		server.Start()
		time.Sleep(1 * time.Millisecond)
		defer func() {
			server.Stop()
		}()

		Convey("Create a new ticket", func() {
			var expectedTitle = "Add killer feature to corcel"
			var expectedBody = "The feature should be..."
			var body = []byte(fmt.Sprintf(`{
				"title" : "%s",
				"body" : "%s"
			}`, expectedTitle, expectedBody))

			req, err := http.NewRequest("POST", server.URL("/tickets/"), bytes.NewBuffer(body))
			req.Close = true
			req.Header.Add("Content-Type", "application/json")

			So(err, ShouldEqual, nil)
			resp, err := client.Do(req)

			So(err, ShouldEqual, nil)

			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			So(resp.Header.Get("Location"), ShouldNotEqual, nil)
			So(resp.Header.Get("Location"), ShouldNotEqual, "")

			Convey("Adding a comment to a ticket", func() {
				var location = resp.Header.Get("Location")

				req, err = http.NewRequest("GET", location, nil)
				req.Close = true
				resp, err = client.Do(req)

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)

				So(err, ShouldEqual, nil)
				var ticket Ticket

				json.Unmarshal(body, &ticket)
				var expectedComment = "Still working on it!"
				body = []byte(fmt.Sprintf(`{
					"comment" : "%s"
				}`, expectedComment))

				req, err = http.NewRequest("POST", server.URL(fmt.Sprintf("/tickets/%s/comments", ticket.ID)), bytes.NewBuffer(body))
				req.Close = true
				req.Header.Add("Content-Type", "application/json")

				So(err, ShouldEqual, nil)

				resp, err = client.Do(req)

				So(err, ShouldEqual, nil)
				So(resp.StatusCode, ShouldEqual, http.StatusCreated)

				Convey("Get a ticket by ID", func() {
					req, err = http.NewRequest("GET", server.URL(fmt.Sprintf("/tickets/%s", ticket.ID)), nil)
					req.Close = true
					resp, err = client.Do(req)

					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)

					So(err, ShouldEqual, nil)

					json.Unmarshal(body, &ticket)
					So(ticket.Status, ShouldEqual, StatusTicketOpen)
					So(len(ticket.Comments), ShouldEqual, 1)
					So(ticket.Comments[0].Comment, ShouldEqual, expectedComment)
				})
			})

			Convey("Closing a ticket", func() {
				var location = resp.Header.Get("Location")

				req, err = http.NewRequest("GET", location, nil)
				req.Close = true
				resp, err = client.Do(req)

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)

				So(err, ShouldEqual, nil)
				var ticket Ticket

				json.Unmarshal(body, &ticket)

				var expectedStatus = StatusTicketComplete
				var expectedComment = "Work complete!"
				body = []byte(fmt.Sprintf(`{
					"status" : %d,
					"comment" : "%s"
				}`, expectedStatus, expectedComment))

				req, err = http.NewRequest("POST", server.URL(fmt.Sprintf("/tickets/%s/close", ticket.ID)), bytes.NewBuffer(body))
				req.Close = true
				req.Header.Add("Content-Type", "application/json")

				So(err, ShouldEqual, nil)

				resp, err = client.Do(req)

				So(err, ShouldEqual, nil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)

				Convey("Get a ticket by ID", func() {
					req, err = http.NewRequest("GET", server.URL(fmt.Sprintf("/tickets/%s", ticket.ID)), nil)
					req.Close = true
					resp, err = client.Do(req)

					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)

					So(err, ShouldEqual, nil)

					json.Unmarshal(body, &ticket)
					So(ticket.Status, ShouldEqual, StatusTicketComplete)
					So(len(ticket.Comments), ShouldEqual, 1)
				})
			})

			Convey("Get a ticket by ID", func() {
				var location = resp.Header.Get("Location")

				req, err = http.NewRequest("GET", location, nil)
				req.Close = true
				resp, err = client.Do(req)

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)

				So(err, ShouldEqual, nil)
				var ticket Ticket

				json.Unmarshal(body, &ticket)

				So(ticket.Title, ShouldEqual, expectedTitle)

				Convey("Update a ticket", func() {
					var updatedTitle = "Something updated in the Title"
					var updatedBody = "Something Updated in Body"
					var body = []byte(fmt.Sprintf(`{
						"title" : "%s",
						"body" : "%s"
					}`, updatedTitle, updatedBody))

					req, err := http.NewRequest("PUT", server.URL("/tickets/"+ticket.ID), bytes.NewBuffer(body))
					req.Close = true
					req.Header.Add("Content-Type", "application/json")
					resp, err = client.Do(req)
					So(err, ShouldEqual, nil)
					So(resp.StatusCode, ShouldEqual, http.StatusOK)
				})
			})

			Convey("Query Multiple Tickets", func() {
				var expectedTitle = "Another title"
				var expectedBody = "Someother Body"
				var body = []byte(fmt.Sprintf(`{
						"title" : "%s",
						"body" : "%s"
					}`, expectedTitle, expectedBody))
				req, _ := http.NewRequest("POST", server.URL("/tickets/"), bytes.NewBuffer(body))
				req.Close = true
				req.Header.Add("Content-Type", "application/json")
				resp, _ := client.Do(req)

				req, _ = http.NewRequest("GET", server.URL("/tickets/"), nil)
				req.Close = true
				req.Header.Add("Content-Type", "application/json")

				resp, _ = client.Do(req)
				defer resp.Body.Close()
				body, _ = ioutil.ReadAll(resp.Body)
				var tickets []Ticket

				json.Unmarshal(body, &tickets)

				So(len(tickets), ShouldEqual, 2)

				Convey("Deleting a Ticket should remove from the query list", func() {

					req, _ = http.NewRequest("DELETE", server.URL(fmt.Sprintf("/tickets/%s", tickets[0].ID)), nil)
					req.Close = true
					resp, _ = client.Do(req)

					So(resp.StatusCode, ShouldEqual, http.StatusOK)
					req, _ = http.NewRequest("GET", server.URL("/tickets/"), nil)
					req.Close = true
					req.Header.Add("Content-Type", "application/json")

					resp, _ = client.Do(req)
					defer resp.Body.Close()
					body, _ = ioutil.ReadAll(resp.Body)

					json.Unmarshal(body, &tickets)

					So(len(tickets), ShouldEqual, 1)
				})
			})
		})

		Convey("Error Handling", func() {
			Convey("Get ticket returns 404 when ticket not found", func() {
				req, _ := http.NewRequest("GET", server.URL(fmt.Sprintf("/tickets/%s", "12345")), nil)
				req.Close = true
				resp, _ := client.Do(req)
				So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
			})

			Convey("Closing a ticket returns 404 when ticket not found", func() {
				var body = []byte("{}")

				var req, err = http.NewRequest("POST", server.URL(fmt.Sprintf("/tickets/%s/close", "12345")), bytes.NewBuffer(body))
				req.Close = true
				req.Header.Add("Content-Type", "application/json")

				So(err, ShouldEqual, nil)

				resp, err := client.Do(req)

				So(err, ShouldEqual, nil)
				So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
			})

			Convey("Update a ticket returns 404 when ticket not found", func() {
				body := []byte("{}")
				req, _ := http.NewRequest("PUT", server.URL("/tickets/"+"12345"), bytes.NewBuffer(body))
				req.Close = true
				req.Header.Add("Content-Type", "application/json")
				resp, _ := client.Do(req)

				So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})

	})
}
