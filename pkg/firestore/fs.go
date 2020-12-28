package firestore

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	fs "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/SunSince90/kube-scraper-backend/pkg/backend"
	"github.com/SunSince90/kube-scraper-backend/pkg/pb"
	"github.com/rs/zerolog"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	timeout = time.Duration(15) * time.Second
)

var (
	log zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type fsBackend struct {
	cache  map[int64]*pb.Chat
	client *fs.Client
	app    *firebase.App
	*Options
	lock sync.Mutex
}

// NewBackend returns a fsHandler, which is an implementation for FS
func NewBackend(ctx context.Context, servAcc string, opts *Options) (backend.Backend, error) {
	// -- Validation
	if len(opts.ChatsCollection) == 0 {
		return nil, fmt.Errorf("no chat collection set")
	}
	if len(opts.ProjectName) == 0 {
		return nil, fmt.Errorf("no project name set")
	}

	// -- Load firebase
	conf := &firebase.Config{ProjectID: opts.ProjectName}
	app, err := firebase.NewApp(ctx, conf, option.WithServiceAccountFile(servAcc))
	if err != nil {
		return nil, err
	}
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	// -- Set up firestore backend
	fs := &fsBackend{
		app:     app,
		client:  fsClient,
		Options: opts,
	}

	if opts.UseCache {
		fs.cache = map[int64]*pb.Chat{}
	}

	return fs, nil
}

// Close the client
func (f *fsBackend) Close() {
	f.client.Close()
}

// GetChatByID gets a chat from firestore
func (f *fsBackend) GetChatByID(id int64) (*pb.Chat, error) {
	// -- Init
	if id == 0 {
		return nil, fmt.Errorf("chat id cannot be 0")
	}

	l := log.With().Str("func", "GetChatByID").Int64("id", id).Logger()
	if f.UseCache {
		// TODO: implement cache
	}

	// -- Get the chat
	docPath := path.Join(f.ChatsCollection, fmt.Sprintf("%d", id))
	ctx, canc := context.WithTimeout(context.Background(), timeout)
	defer canc()

	doc, err := f.client.Doc(docPath).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, backend.ErrNotFound
		}

		return nil, err
	}
	l.Debug().Msg("pulled from firestore")

	// -- Cast and return
	var _chat chat
	if err := doc.DataTo(&_chat); err != nil {
		return nil, err
	}
	c := convertToProto(&_chat)

	if f.UseCache {
		// TODO: implement cache
	}

	return c, nil
}

// GetChatByUsername gets a chat from firestore by username
func (f *fsBackend) GetChatByUsername(username string) (*pb.Chat, error) {
	// TODO: implement me
	return nil, nil
}

// StoreChats inserts a chat into firestore
func (f *fsBackend) StoreChat(c *pb.Chat) error {
	// TODO: implement me
	return nil
}

// DeleteChat deletes a chat from firestore
func (f *fsBackend) DeleteChat(id int64) error {
	// TODO: implement me
	return nil
}

// GetAllChats gets all chat from firestore
func (f *fsBackend) GetAllChats() ([]*pb.Chat, error) {
	// TODO: implement me
	return nil, nil
}