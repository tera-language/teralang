package main

import (
	"encoding/json"
	"fmt"

	"github.com/tera-language/teralang/internal/logger"
	"github.com/tera-language/teralang/internal/parser"
)

func main() {
	program, err := parser.Parse("./testdata/test.tera")
	if err != nil {
		logger.Error(err)
	}
	prettyJson, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		logger.Error(err)
	}
	fmt.Printf("%s\n", string(prettyJson))
}
