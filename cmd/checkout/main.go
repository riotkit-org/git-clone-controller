package serve

type Command struct {
	LogLevel string
	Path     string
	Url      string
	Username string
	Token    string
	Revision string
}

func (c *Command) Run() error {
	return nil
}
