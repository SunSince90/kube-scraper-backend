package firestore

type firestoreOptions struct {
	chatsCollection    string
	debug              bool
	listen             bool
	serviceAccountPath string
	projectName        string
	address            string
	port               int
}
