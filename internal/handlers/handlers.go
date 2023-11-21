package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"temlate/proto/proto/pb"
)

type service interface {
	Create(msg string, prefid string) (string, error)
	Get(id string) (string, error)
	GetAll() (string, error)
}

type Handlers struct {
	service
	pb.UnimplementedGatewayTemplateServer
}

func New(service service) *Handlers {
	return &Handlers{service: service, UnimplementedGatewayTemplateServer: pb.UnimplementedGatewayTemplateServer{}}
}
func (h *Handlers) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	// Read the param noteId
	//link := c.Params("Link")
	fmt.Println("Post - ")
	//var shortUrl model.Links
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	type Template struct {
		Prefid string `json:"preferid"`
		Temp   string `json:"template"`
	}
	var Temp Template = Template{Temp: request.Template, Prefid: request.Preferid}
	fmt.Println("Post - ", Temp.Temp, Temp.Prefid)
	msg, err := h.service.Create(Temp.Temp, Temp.Prefid)
	if err != nil {
		log.Error(err.Error() + msg)
		return &pb.CreateResponse{Template: http.StatusText(400)}, err
	}
	return &pb.CreateResponse{Template: http.StatusText(200)}, nil
	// Find the note with the given id
	// Return the note with the id
}

func (h *Handlers) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	id := in.Id
	fmt.Println("Get - ", id)
	msg, err := h.service.Get(id)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{Template: msg}, nil
}

func (h *Handlers) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	fmt.Println("Get all- ")
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	msg, err := h.service.GetAll()
	if err != nil {
		log.Error(msg)
		return &pb.GetAllResponse{Template: strconv.Itoa(http.StatusNotFound)}, err
	}
	return &pb.GetAllResponse{Template: msg}, nil
}
