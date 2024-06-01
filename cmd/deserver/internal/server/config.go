package server

// Config contains all the necessary configuration options for the server.
type Config struct {
	LisAddr  string
	HTTPPort string
}

var defaultConfig = Config{
	LisAddr:  "localhost:3030",
	HTTPPort: ":8081",
}
