package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
	"temlate/config"
	"temlate/internal/handlers"
	"temlate/internal/repository"
	"temlate/internal/service"
)

type Repository interface {
	CreateTemplate(msg string) (string, error)
	GetTemplate(id int) (string, error)
	GetTemplates() (string, error)
}

type App struct {
	handlers   *handlers.Handlers
	service    *service.Service
	repository *repository.Repository
	routers    *fiber.App
}

func New(ctx context.Context, cfg config.Config) *App {
	fmt.Println("Start")
	//logger := slog.Logger{} //засунуть в контекст
	app := &App{}
	app.routers = fiber.New()
	app.repository = repository.New(cfg)
	app.service = service.New(app.repository)
	app.handlers = handlers.New(app.service)
	return app
}

func Run(app *App) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app.routers.Post("/", app.handlers.Post)
	app.routers.Get("/:id", app.handlers.Get)
	app.routers.Get("/", app.handlers.GetAll)
	err := app.routers.Listen(":3000")
	if err != nil {
		log.Error(err.Error())
		return
	}
}
