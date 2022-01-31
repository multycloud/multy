package cli

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func StartCli() {
	commands := []Command{&TranslateCommand{}, &ApplyCommand{}, &DestroyCommand{}}

	flagset := flag.NewFlagSet("cmd", flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)
	_ = flagset.Parse(os.Args[1:])

	args := flagset.Args()
	if len(args) == 0 {
		log.Fatalf("no command was specified")
	}

	var selected Command
	for _, c := range commands {
		if c.Name() == args[0] {
			selected = c
			break
		}
	}

	if selected == nil {
		log.Fatalf("command not found: %s", args[0])
	}

	selected.Init()
	flag.Parse()
	args = flag.Args()

	err := selected.Execute(args[1:])
	if err != nil {
		log.Fatalf(err.Error())
	}

}
