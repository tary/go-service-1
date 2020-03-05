package mdns

import (
	"context"

	"github.com/giant-tech/go-service/base/plugin/registry"
)

type portType struct{}

// Port if it is not 0, replace the port 5353 with this port number.
func Port(port int) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, portType{}, port)
	}
}
