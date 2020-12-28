package firestore

import (
	"github.com/SunSince90/kube-scraper-backend/pkg/pb"
)

func convertToProto(c *chat) *pb.Chat {
	return &pb.Chat{
		Id:        c.ChatID,
		Type:      c.Type,
		Username:  c.Username,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
}
