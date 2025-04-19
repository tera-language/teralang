package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/tera-language/teralang/internal/logger"
	"github.com/tera-language/teralang/internal/parser"
	"github.com/tera-language/teralang/internal/server"
)

func main() {
	help := false
	flag.BoolVar(&help, "help", false, "Display this help message and exit.")
	flag.BoolVar(&help, "h", false, "Display this help message and exit.")

	port := 3000
	flag.IntVar(&port, "port", 3000, "Set the port number.")
	flag.IntVar(&port, "p", 3000, "Set the port number.")

	flag.Parse()

	if help {
		fmt.Print(`
USAGE: teralang <path> [-p <port>]

ARGUMENTS:
  <path>               The path to the entrypoint .tera file.

OPTIONS:
  --port, -p <port>    Set the port number.
  --help, -h           Display this help message and exit.
`)
		return
	}

	entrypoint := flag.Arg(0)
	logger.Infoln("Starting parsing...")
	program, err := parser.Parse(entrypoint)
	if err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}
	logger.Successln("Parsing done!")

	mux := server.Server(program)
	logger.Successf("Server started at http://localhost:%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}
}
