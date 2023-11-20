package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
	"os"
)

type service interface {
	Create(msg string, prefid string) (string, error)
	Get(id string) (string, error)
	GetAll() (string, error)
}

type Handlers struct {
	service
}

func New(service service) *Handlers {
	return &Handlers{service}
}
func (h *Handlers) Post(c *fiber.Ctx) error {
	// Read the param noteId
	//link := c.Params("Link")
	//var shortUrl model.Links
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	type Template struct {
		Prefid string `json:"preferid"`
		Temp   string `json:"template"`
	}
	var Temp Template
	err := c.BodyParser(&Temp)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	fmt.Println("Post - ", Temp.Temp, Temp.Prefid)
	msg, err := h.service.Create(Temp.Temp, Temp.Prefid)
	if err != nil {
		log.Error(err.Error() + msg)
		return c.SendStatus(http.StatusBadRequest)
	}
	return c.SendStatus(http.StatusOK)
	// Find the note with the given id
	// Return the note with the id
}

func (h *Handlers) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println("Get - ", id)
	msg, err := h.service.Get(id)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	return c.SendString(msg)
}

func (h *Handlers) GetAll(c *fiber.Ctx) error {
	fmt.Println("Get all- ")
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	msg, err := h.service.GetAll()
	if err != nil {
		log.Error(msg)
		return c.SendStatus(http.StatusNotFound)
	}
	return c.SendString(msg)
}
