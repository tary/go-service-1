package app

import (
	"github.com/giant-tech/go-service/framework/app/internal"
)

// interface cmd接口

type Cmd interface {

	// The cli app within this cmd
	App() *internal.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
}

// NewCmd new
func NewCmd(opts ...Option) Cmd {
	return newCmd(opts...)
}
