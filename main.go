package main

import (
	"flag"
	"fmt"

	"github.com/tera-language/teralang/internal/logger"
	"github.com/tera-language/teralang/internal/parser"
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
}
