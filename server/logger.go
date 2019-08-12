package server

import "github.com/go-kit/kit/log"

//Logger fazzkit logger option
type Logger struct {
	Logger    log.Logger
	Namespace string
	Subsystem string
	Action    string
}
