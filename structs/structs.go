package structs

// Config is a struct to hold the config parameters provided on startup
type Config struct {
	Port    int16
	Tickets string
	Users   string
	Cert    string
	Key     string
}
