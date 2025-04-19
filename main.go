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

	flag.Parse()

	if help {
		fmt.Print(`
USAGE: teralang <path>

ARGUMENTS:
  <path>        The path to the entrypoint .tera file.

OPTIONS:
  --help, -h    Display this help message and exit.
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
	logger.Successln("Server started at http://localhost:3000")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}
}
