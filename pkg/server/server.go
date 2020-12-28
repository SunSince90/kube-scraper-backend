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

func (b *backendServer) GetChat(ctx context.Context, r *pb.ChatRequest) (*pb.ChatResponse, error) {
	// TODO: implement me
	return nil, nil
}

func (b *backendServer) GetAllChats(ctx context.Context, r *pb.ChatRequest) (*pb.ChatResponse, error) {
	// TODO: implement me
	return nil, nil
}
