package app

// Name of the service
func Name(n string) OptionFunc {
	return func(o *Options) {

	}
}

// Version of the service
func Version(v string) OptionFunc {
	return func(o *Options) {

	}
}

// Before and Afters

// BeforeStart run funcs before service starts
func BeforeStart(fn func() error) OptionFunc {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func() error) OptionFunc {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart run funcs after service starts
func AfterStart(fn func() error) OptionFunc {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func() error) OptionFunc {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
