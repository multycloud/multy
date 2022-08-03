package main

import (
	"github.com/multycloud/multy/cmd"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		// pprof
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	cmd.StartCli()
}
