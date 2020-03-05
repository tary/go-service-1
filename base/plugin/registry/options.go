package registry

import (
	"context"
	"time"
)

// Options options
type Options struct {
	Addrs   []string
	Timeout time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RegisterOptions register options
type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// WatchOptions watch options
type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Timeout timeout
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// RegisterTTL ttl
func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// WatchService Watch a service
/*func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}
*/
