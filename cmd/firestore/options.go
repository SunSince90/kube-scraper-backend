package firestore

type firestoreOptions struct {
	chatsCollection    string
	debug              bool
	listen             bool
	serviceAccountPath string
	projectID          string
	address            string
	port               int
}
