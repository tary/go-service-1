package app

import (
	"context"

	"github.com/giant-tech/go-service/base/itf/ioption"
	"github.com/giant-tech/go-service/base/plugin/registry"
)

// Options options
type Options struct {
	// For the Command Line itself
	ID          string
	Name        string
	Description string
	Version     string

	// We need pointers to things so we can swap them out if needed.
	//Logger *ilog.Logger
	//Registry  *registry.Registry
	//Selector  *selector.Selector
	//Connector *connector.Connector
	//Transport *transport.Transport
	//Client    *client.Client
	//Server    *server.Server

	//Loggers    map[string]func(...log.Option) log.Logger
	//Clients    map[string]func(...client.Option) client.Client
	Registries map[string]func(...registry.Option) registry.IRegistry
	//Selectors  map[string]func(...selector.Option) selector.Selector
	//Connectors map[string]func(...connector.Option) connector.Connector
	//Servers    map[string]func(...server.Option) server.Server
	//	Transports map[string]func(...transport.Option) transport.Transport

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option option
type Option func(*Options)

// Name of the service
func Name(n string) ioption.OptionFunc {
	return func(o *ioption.Options) {

	}
}

// Version of the service
func Version(v string) ioption.OptionFunc {
	return func(o *ioption.Options) {

	}
}

// Before and Afters

// BeforeStart run funcs before service starts
func BeforeStart(fn func() error) ioption.OptionFunc {
	return func(o *ioption.Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func() error) ioption.OptionFunc {
	return func(o *ioption.Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart run funcs after service starts
func AfterStart(fn func() error) ioption.OptionFunc {
	return func(o *ioption.Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func() error) ioption.OptionFunc {
	return func(o *ioption.Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
