package cmd

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"time"
)

func StartCli() {
	helpCmd := &HelpCommand{}
	commands := []Command{&VersionCommand{}, &ServeCommand{}, &ListCommand{}, &DeleteCommand{}, helpCmd}
	helpCmd.AvailableCommands = commands

	if len(os.Args) < 2 {
		err := helpCmd.Execute(nil)
		if err != nil {
			log.Fatalf("unable to show help command: %s", err.Error())

		}
		return
	}

	var selected Command
	for _, c := range commands {
		if c.Description().Name == os.Args[1] {
			selected = c
			break
		}
	}

	if selected == nil {
		log.Fatalf("command not found: %s", os.Args[1])
	}

	f := flag.NewFlagSet(selected.Description().Name, flag.ExitOnError)
	f.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", selected.Description().Name, selected.Description().Description)
		if selected.Description().Usage != "" {
			_, _ = fmt.Fprintf(os.Stderr, "Usage: %s\n", selected.Description().Usage)
		}
		f.PrintDefaults()
	}

	selected.ParseFlags(f, os.Args[2:])

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

	err := selected.Execute(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
