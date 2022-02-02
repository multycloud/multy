package cli

import "context"

type Command interface {
	Name() string
	Init()
	Execute(args []string, ctx context.Context) error
}
