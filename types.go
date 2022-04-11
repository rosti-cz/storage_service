package main

// Message coming from the admin. Message is coming from the admin interface and
// it says that something happening there and we should check if we should do something with it.
type Message struct {
	EventType  string   `json:"event_type"`
	DBID       int      `json:"db_id"`
	DBName     string   `json:"db_name"`
	Username   string   `json:"username"`
	UsernameRO string   `json:"username_ro"`
	Password   string   `json:"password"`
	PasswordRO string   `json:"password_ro"`
	Extensions []string `json:"extensions"`
}

// State is async response back to the admin and it says if something was done.
// If admin tells us that database was created we should create it in the database instance locally
// and return State message back to admin.
type State struct {
	DBID    int    `json:"db_id"`
	DBName  string `json:"db_name"`
	Error   bool   `json:"error"`   // true if there was an error
	Message string `json:"message"` // error message or state like created,password_changed or deleted
}

// Backend is interface to handle databases
type Backend interface {
	CreateUser(user, password, database string) error
	CreateROUser(user, password, database string) error
	CreateDatabase(database, owner string, extensions []string) error
	ChangePassword(user, password string) error
	DropUser(user string) error
	DropDatabase(database string) error
}

// Metrics is used to share status of the service with the ecosystem
type Metrics struct {
	Messages int    `json:"messages"`
	Service  string `json:"service"`
}

// Structure for metrics message compatible with metrics-receiver
type ReceiverMetric struct {
	Lines   []string `json:"lines"` // prometheus like metrics, one per line
	Service string   `json:"service"`
}
