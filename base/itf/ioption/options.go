package ioption

import "context"

// Options options
type Options struct {
	//Cmd       cmd.Cmd

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Signal bool
}

func newOptions(opts ...OptionFunc) Options {
	opt := Options{
		Context: context.Background(),
		Signal:  true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// OptionFunc option函数
type OptionFunc func(*Options)
