package cli

type Command interface {
	Name() string
	Init()
	Execute(args []string) error
}
