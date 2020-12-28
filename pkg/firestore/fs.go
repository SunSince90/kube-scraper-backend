package firestore

import (
	"context"
	"errors"
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
	"google.golang.org/api/iterator"
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
	// -- Init
	if len(username) == 0 {
		return nil, fmt.Errorf("chat username cannot be 0")
	}

	l := log.With().Str("func", "GetChatByUsername").Str("username", username).Logger()
	ctx, canc := context.WithTimeout(context.Background(), timeout)
	defer canc()

	// -- Get the chat
	docIter := f.client.Collection(f.ChatsCollection).Where("username", "==", username).Limit(1).Documents(ctx)
	doc, err := docIter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
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

// StoreChats inserts a chat into firestore
func (f *fsBackend) StoreChat(c *pb.Chat) error {
	// -- Init
	if c == nil {
		return fmt.Errorf("chat cannot be nil")
	}
	l := log.With().Str("func", "StoreChat").Int64("id", c.Id).Logger()
	addChat := chat{
		ChatID:    c.Id,
		Type:      c.Type,
		Username:  c.Username,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}

	// -- Store the chat on firestore
	docPath := path.Join(f.ChatsCollection, fmt.Sprintf("%d", c.Id))
	ctx, canc := context.WithTimeout(context.Background(), timeout)
	defer canc()

	_, err := f.client.Doc(docPath).Set(ctx, addChat)
	if err != nil {
		return err
	}

	if f.UseCache {
		// TODO: implement cache
	}

	l.Debug().Msg("stored on firestore")
	return nil
}

// DeleteChat deletes a chat from firestore
func (f *fsBackend) DeleteChat(id int64) error {
	// -- Init
	if id == 0 {
		return fmt.Errorf("chat id cannot be 0")
	}

	// -- Delete and return
	docPath := path.Join(f.ChatsCollection, fmt.Sprintf("%d", id))
	ctx, canc := context.WithTimeout(context.Background(), timeout)
	defer canc()

	_, err := f.client.Doc(docPath).Delete(ctx)

	if f.UseCache {
		// TODO: implement cache
	}

	return err
}

// GetAllChats gets all chat from firestore
func (f *fsBackend) GetAllChats() ([]*pb.Chat, error) {
	// -- Init
	l := log.With().Str("func", "GetAllChats").Logger()
	if f.UseCache {
		// TODO: implement cache
	}

	ctx, canc := context.WithTimeout(context.Background(), timeout)
	defer canc()

	// -- Get list
	list := []*pb.Chat{}
	dociter := f.client.Collection(f.ChatsCollection).Documents(ctx)
	defer dociter.Stop()

	for {
		doc, err := dociter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}

			return nil, err
		}

		var _chat chat
		if err := doc.DataTo(&_chat); err != nil {
			l.Err(err).Int64("id", _chat.ChatID).Msg("error while trying to get this document, skipping...")
			continue
		}

		c := convertToProto(&_chat)

		if f.UseCache {
			// TODO: implement cache
		}

		list = append(list, c)
	}

	return list, nil
}
