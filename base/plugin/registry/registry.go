// Package registry is an interface for service discovery
package registry

import "errors"

// Registry The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

// Option option
type Option func(*Options)

// RegisterOption register option
type RegisterOption func(*RegisterOptions)

// WatchOption watch option
type WatchOption func(*WatchOptions)

var (
	// ErrNotFound Not found error when GetService is called
	ErrNotFound = errors.New("not found")
	// ErrWatcherStopped Watcher stopped error when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)
