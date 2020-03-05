package app

// interface cmd接口

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

func newCmd(opts ...Option) Cmd {

	options := Options{}
	cmd := new(cmd)
	cmd.opts = options

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
