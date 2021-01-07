package server

import (
	"context"
	"fmt"
	"os"

	"github.com/SunSince90/kube-scraper-backend/pkg/backend"
	"github.com/SunSince90/kube-scraper-backend/pkg/pb"
	"github.com/rs/zerolog"
)

var (
	log zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type backendServer struct {
	backend backend.Backend
	pb.UnimplementedBackendServer
}

// New returns a new instance of the server using the provided backend
func New(b backend.Backend) (pb.BackendServer, error) {
	if b == nil {
		return nil, fmt.Errorf("no backend set")
	}

	serv := &backendServer{
		backend: b,
	}

	return serv, nil
}

// GetChat gets chat from the backend
func (b *backendServer) GetChat(ctx context.Context, r *pb.ChatRequest) (*pb.ChatResponse, error) {
	var chat *pb.Chat
	var err error

	// -- Get the chat by ID
	if r.Id != 0 {
		chat, err = b.backend.GetChatByID(r.Id)
	}

	// -- Get the chat by username
	if len(r.Username) > 0 {
		chat, err = b.backend.GetChatByUsername(r.Username)
	}

	if err != nil {
		return &pb.ChatResponse{
			Code:    500,
			Message: err.Error(),
		}, fmt.Errorf("error while getting chat")
	}

	return &pb.ChatResponse{
		Code:    200,
		Message: "ok",
		Chats:   []*pb.Chat{chat},
	}, nil
}

// GetChatsList returns a list of chats according to the request
func (b *backendServer) GetChatsList(ctx context.Context, r *pb.ChatRequest) (*pb.ChatResponse, error) {
	chats, err := b.backend.GetAllChats()
	if err != nil {
		return &pb.ChatResponse{
			Code:    500,
			Message: err.Error(),
		}, fmt.Errorf("error while getting chats list")
	}

	// Has a filter?
	if len(r.Type) == 0 {
		return &pb.ChatResponse{
			Code:    200,
			Message: "ok",
			Chats:   chats,
		}, nil
	}

	if r.Type != "group" && r.Type != "supergroup" && r.Type != "private" || r.Type != "channel" {
		return &pb.ChatResponse{
			Code:    500,
			Message: "unrecognized chat type",
			Chats:   chats,
		}, fmt.Errorf("unrecognized chat type (%s)", r.Type)
	}

	_chats := []*pb.Chat{}
	for _, chat := range chats {
		if chat.Type == r.Type {
			_chats = append(_chats, chat)
		}
	}

	return &pb.ChatResponse{
		Code:    200,
		Message: "ok",
		Chats:   _chats,
	}, nil
}
