package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/braintree/manners"
	"github.com/nu7hatch/gouuid"
	"github.com/uber-go/zap"
	"gopkg.in/gin-gonic/gin.v1"
)

//TicketServer ...
type TicketServer struct {
	Host              string
	Port              int
	Server            *manners.GracefulServer
	Store             TicketStore
	CommandController CommandController
	Engine            *gin.Engine
}

//CreateTicketServer ...
func CreateTicketServer(port int, host string) *TicketServer {
	var engine = gin.Default()

	var httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: engine,
		//TODO: CONFIG
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var server = &TicketServer{
		Host:   host,
		Port:   port,
		Server: manners.NewWithServer(httpServer),
		Store: &LoggingTicketStore{
			WrappedTicketStore: &InProcessTicketStore{
				tickets: []Ticket{},
			},
		},
		CommandController: CommandController{
			HandlerMappings: map[string]CommandHandler{},
		},
		Engine: engine,
	}

	server.setupCommandHandlers()

	server.setupRouting()

	logger.Info("TicketServer Created",
		zap.String("host", host),
		zap.Int("port", port),
	)

	return server
}

//URL ...
func (instance *TicketServer) URL(path string) string {
	return fmt.Sprintf("http://%s:%d%s", instance.Host, instance.Port, path)
}

func (instance *TicketServer) setupCommandHandlers() {
	instance.CommandController.AddHandler(NameCreateTicketCommand, CreateTicketCommandHandler{store: instance.Store})
	instance.CommandController.AddHandler(NameDeleteTicketCommand, DeleteTicketCommandHandler{store: instance.Store})
	instance.CommandController.AddHandler(NameUpdateTicketCommand, UpdateTicketCommandHandler{store: instance.Store})
	instance.CommandController.AddHandler(NameCloseTicketCommand, CloseTicketCommandHandler{store: instance.Store})
	instance.CommandController.AddHandler(NameAddTicketCommentCommand, AddTicketCommentCommandHandler{store: instance.Store})

	logger.Info("Command Handlers Setup")
}

func (instance *TicketServer) setupRouting() {
	tickets := instance.Engine.Group("/tickets")
	{
		tickets.GET("/", func(c *gin.Context) {
			c.JSON(200, instance.Store.Query())
		})
		tickets.GET("/:id", func(c *gin.Context) {
			var id = c.Param("id")
			var ticket, err = instance.Store.GetByID(id)
			if err != nil {
				logger.Error(fmt.Sprintf("%v", err),
					zap.String("id", id),
				)
				c.Status(http.StatusNotFound)
			} else {
				c.JSON(200, ticket)
			}
		})
		tickets.POST("/", func(c *gin.Context) {
			var command CreateTicketCommand
			if c.BindJSON(&command) == nil {
				id, _ := uuid.NewV4()
				command.ID = id.String()
				err := instance.CommandController.Handle(NameCreateTicketCommand, command)
				if err != nil {
					logger.Error(fmt.Sprintf("%v", err))
					c.Status(http.StatusInternalServerError)
				} else {
					c.Header("Location", fmt.Sprintf("http://%s:%d/tickets/%s", instance.Host, instance.Port, id))
					c.Status(http.StatusCreated)
				}
			}
		})
		tickets.POST("/:id/close", func(c *gin.Context) {
			var command CloseTicketCommand
			if c.BindJSON(&command) == nil {
				var id = c.Param("id")
				command.ID = id
				logger.Debug("Closing Ticket",
					zap.String("id", command.ID),
					zap.Int("status", command.Status),
					zap.String("comment", command.Comment),
				)
				err := instance.CommandController.Handle(NameCloseTicketCommand, command)
				if err != nil {
					logger.Error(fmt.Sprintf("%v", err),
						zap.String("id", id),
					)
					c.Status(http.StatusNotFound)
				} else {
					c.Status(http.StatusOK)
				}
			}
		})
		tickets.POST("/:id/comments", func(c *gin.Context) {
			var command AddTicketCommentCommand
			if c.BindJSON(&command) == nil {
				var id = c.Param("id")
				command.ID = id
				err := instance.CommandController.Handle(NameAddTicketCommentCommand, command)
				if err != nil {
					logger.Error(fmt.Sprintf("%v", err),
						zap.String("id", id),
					)
					c.Status(http.StatusInternalServerError)
				} else {
					c.Status(http.StatusCreated)
				}
			}
		})
		tickets.PUT("/:id", func(c *gin.Context) {
			var command UpdateTicketCommand
			if c.BindJSON(&command) == nil {
				var id = c.Param("id")
				command.ID = id
				err := instance.CommandController.Handle(NameUpdateTicketCommand, command)
				if err != nil {
					logger.Error(fmt.Sprintf("%v", err),
						zap.String("id", id),
					)
					c.Status(http.StatusNotFound)
				} else {
					c.Header("Location", fmt.Sprintf("http://%s:%d/tickets/%s", instance.Host, instance.Port, c.Param("id")))
					c.JSON(http.StatusOK, nil)
				}
			}
		})
		tickets.DELETE("/:id", func(c *gin.Context) {
			var id = c.Param("id")
			var command = DeleteTicketCommand{
				ID: id,
			}
			err := instance.CommandController.Handle(NameDeleteTicketCommand, command)
			if err != nil {
				logger.Error(fmt.Sprintf("%v", err),
					zap.String("id", id),
				)
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(http.StatusOK)
			}
		})
	}
	logger.Info("Routing Setup")
}

//Start ...
func (instance *TicketServer) Start() {
	go instance.Server.ListenAndServe()
	logger.Info("TicketServer Started")
}

//Stop ...
func (instance *TicketServer) Stop() {
	instance.Server.BlockingClose()
	logger.Info("TicketServer Stopped")
}

var (
	logger zap.Logger
)

func init() {
	InitializeLogger()
}

//InitializeLogger ...
func InitializeLogger() {
	logger = zap.New(
		zap.NewJSONEncoder(), // drop timestamps in tests
		zap.DebugLevel,
	)
}

//SetLogger ...
func SetLogger(newLogger zap.Logger) {
	logger = newLogger
}

func main() {

	var port = 45000
	var host = "localhost"

	var server = CreateTicketServer(port, host)
	server.Start()
}
