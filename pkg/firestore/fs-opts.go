package firestore

// Options contains options for firestore
type Options struct {
	// ProjectID is the id of the firebase project
	ProjectID string `yaml:"projectId"`
	// ChatsCollection is the name of the collections where
	// chats are stored
	ChatsCollection string `yaml:"chatsCollection"`
	// UseCache tells whether to cache chats locally
	UseCache bool `yaml:"useCache"`
}
