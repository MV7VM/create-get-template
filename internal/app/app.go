package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	"os"
	"temlate/config"
	"temlate/internal/handlers"
	"temlate/internal/repository"
	"temlate/internal/service"
	"temlate/proto/proto/pb"
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
}

const (
	grpcPort = ":8090"
	httpPort = ":3000"
)

func New(ctx context.Context, cfg config.Config) *App {
	fmt.Println("Start")
	//logger := slog.Logger{} //засунуть в контекст
	app := &App{}
	app.repository = repository.New(cfg)
	app.service = service.New(app.repository)
	app.handlers = handlers.New(app.service)
	return app
}

func Run(app *App) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	lis, err := net.Listen("tcp", grpcPort)
	s := grpc.NewServer()
	pb.RegisterGatewayTemplateServer(s, app.handlers)
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Error("failed to serve: " + err.Error())
			//log.Error("failed to serve: " + err.Error())
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		grpcPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to dial serve: " + err.Error())
		//log.Error("Failed to dial server: " + err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection " + err.Error())
		}
	}(conn)

	gwmux := runtime.NewServeMux()
	err = pb.RegisterGatewayTemplateHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Error("failed : " + err.Error())
		log.Error("Failed to register gateway:" + err.Error())
	}

	gwServer := &http.Server{
		Addr:    httpPort,
		Handler: gwmux,
	}

	log.Info("Serving gRPC-Gateway on port " + httpPort)
	if err = gwServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Warn("Server closed: " + err.Error())
			os.Exit(0)
		}
		log.Error("Failed to listen and serve: " + err.Error())
	}
	if err != nil {
		return
	}
}
