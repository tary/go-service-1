package app

import "github.com/giant-tech/go-service/base/itf/ioption"

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
