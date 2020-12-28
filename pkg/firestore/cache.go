package firestore

import "github.com/SunSince90/kube-scraper-backend/pkg/pb"

func (f *fsBackend) getChatFromCache(id int64) *pb.Chat {
	f.lock.Lock()
	defer f.lock.Unlock()

	c, exists := f.cache[id]
	if exists && c != nil {
		return c
	}

	return nil
}

func (f *fsBackend) insertChatIntoCache(c *pb.Chat) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.cache[c.Id] = c
}

func (f *fsBackend) deleteChatFromCache(id int64) {
	f.lock.Lock()
	defer f.lock.Unlock()

	delete(f.cache, id)
}

func (f *fsBackend) getAllChatsFromCache() []*pb.Chat {
	f.lock.Lock()
	defer f.lock.Unlock()

	if len(f.cache) == 0 {
		return []*pb.Chat{}
	}

	list := make([]*pb.Chat, len(f.cache))
	i := 0

	for _, c := range f.cache {
		list[i] = c
		i++
	}

	return list
}
