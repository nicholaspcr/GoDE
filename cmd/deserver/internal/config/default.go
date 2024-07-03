package config

import "github.com/nicholaspcr/GoDE/internal/log"

func Default() *DeServer {
	return &DeServer{
		Log: log.Config,
	}
}
