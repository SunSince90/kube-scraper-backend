package backend

import (
	"context"

	"github.com/SunSince90/kube-scraper-backend/pkg/pb"
)

// Backend is a backend that can be used to store and retrieve
// chats.
type Backend interface {
	// ListenForChats listens for chats and updates the cache accordingly
	ListenForChats(context.Context, chan struct{})
	// GetChatByID retrieves a chat by the id
	GetChatByID(int64) (*pb.Chat, error)
	// GetChatByUsername retrieves a chat by username
	GetChatByUsername(string) (*pb.Chat, error)
	// GetAllChatIDs gets all chats
	GetAllChats() ([]*pb.Chat, error)
	// StoreChat stores a new chat in the database
	StoreChat(*pb.Chat) error
	// DeleteChat deletes a chat from the database
	DeleteChat(int64) error
	// Close any client
	Close()
}
