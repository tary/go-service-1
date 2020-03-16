package app

import (
	"github.com/giant-tech/go-service/base/plugin/registry"
	"github.com/giant-tech/go-service/base/plugin/registry/mdns"
)

// Cmd interface cmd接口
type Cmd interface {

	// The cli app within this cmd
	//App() *internal.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
}

type cmd struct {
	opts Options
	//app  *cli.App
}

var (
	// DefaultRegistries default registries
	DefaultRegistries = map[string]func(...registry.Option) registry.IRegistry{
		//"consul": consul.NewRegistry,
	}

	defaultRegistry = "consul"
)

func newRegistry(name string) registry.IRegistry {
	var r registry.IRegistry
	if name == "mdns" {
		// Many products will also use the 5353 port for mdns service discovery, such as jenkins.
		// Change to other ports, to avoid 'discovering other APP services,
		// and protocol parsing failed, resulting in an infinite loop BUG'
		opt := []registry.Option{mdns.Port(5354)}
		r = DefaultRegistries[name](opt...)
	} else {
		r = DefaultRegistries[name]()
	}
	return r
}

func newCmd(opts ...Option) Cmd {

	//l := newLogger(defaultLog)
	/*r := */
	newRegistry(defaultRegistry)
	/*tran := newTransport(defaultTransport)
	slt := newSelector(defaultSelector, []selector.Option{
		selector.Registry(r),
	}...)
	ct := newConnctor(defaultConnector, []connector.Option{
		connector.Transport(tran),
	}...)
	srv := newServer(defaultServer, []server.Option{
		server.Registry(r),
		server.Transport(tran),
	}...)
	c := newClient(defaultClient, []client.Option{
		client.Registry(r),
		client.Transport(tran),
		client.Selector(slt),
		client.Connector(ct),
	}...)
	*/

	options := Options{
		/*Logger:     &l,
		Client:     &c,
		Registry:   &r,
		Server:     &srv,
		Selector:   &slt,
		Connector:  &ct,
		Transport:  &tran,
		Loggers:    DefaultLogs,
		Clients:    DefaultClients,
		Registries: DefaultRegistries,
		Selectors:  DefaultSelectors,
		Servers:    DefaultServers,
		Transports: DefaultTransports,
		Action:     func(c *cli.Context) {},
		*/
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Description) == 0 {
		options.Description = "a v-micro service"
	}

	cmd := new(cmd)
	cmd.opts = options
	/*cmd.app = cli.NewApp()
	cmd.app.Name = cmd.opts.Name
	cmd.app.Version = cmd.opts.Version
	cmd.app.Usage = cmd.opts.Description
	cmd.app.Before = cmd.Before
	cmd.app.Flags = DefaultFlags
	cmd.app.Action = func(c *cli.Context) {}

	if len(options.Version) == 0 {
		cmd.app.HideVersion = true
	}
	*/
	return cmd
}

/*func (c *cmd) App() *cli.App {
	return c.app
}
*/

func (c *cmd) Options() Options {
	return c.opts
}

func (c *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	/*c.app.Name = c.opts.Name
	c.app.Version = c.opts.Version
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description
	c.app.Flags = append(c.app.Flags, c.opts.Flags...)
	c.app.Action = c.opts.Action

	if err := c.app.Run(os.Args); err != nil {
		log.Fatal(err)
	}*/

	return nil
}

// NewCmd new
func NewCmd(opts ...Option) Cmd {
	return newCmd(opts...)
}
