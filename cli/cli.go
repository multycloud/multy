package cli

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

func StartCli() {
	commands := []Command{&TranslateCommand{}, &ApplyCommand{}, &DestroyCommand{}, &VersionCommand{}}

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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			fmt.Println("Cancel signal received. Terminating...")
			cancel()
			go func() {
				// If it hasn't stopped gracefully by now, just bruteforce it.
				time.Sleep(5 * time.Second)
				os.Exit(1)
			}()
		case <-ctx.Done():
		}
	}()

	err := selected.Execute(args[1:], ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}

}
