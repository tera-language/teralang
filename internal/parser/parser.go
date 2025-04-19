package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/tera-language/teralang/internal/logger"
	tree_sitter_teralang "github.com/tera-language/tree-sitter-teralang/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

var parsedPaths []string = []string{}

func Parse(filepath string) ([]route, error) {
	parser := tree_sitter.NewParser()
	lang := tree_sitter.NewLanguage(tree_sitter_teralang.Language())

	err := parser.SetLanguage(lang)
	if err != nil {
		return nil, err
	}

	text, err := os.ReadFile(path.Join(".", filepath))
	if err != nil {
		return nil, err
	}

	tree := parser.Parse(text, nil)
	rootNode := tree.RootNode()

	filepath = path.Clean(filepath)
	parsedPaths = append(parsedPaths, filepath)

	program := []route{}
	program, err = parseNode(rootNode, text, filepath, program)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func parseNode(node *tree_sitter.Node, source []byte, sourcePath string, program []route) ([]route, error) {
	switch node.GrammarName() {
	case "import":
		start, end := node.NamedChild(0).ByteRange()
		importPath := string(source[start+1 : end-1]) // skip "" characters
		relativePath := path.Join(path.Dir(sourcePath), importPath)
		if slices.Contains(parsedPaths, relativePath) {
			logger.Warning("Already parsed: ", relativePath)
			break
		}
		importedProgram, err := Parse(relativePath)
		if err != nil {
			return nil, err
		}
		program = importedProgram
	case "route":
		start, end := node.NamedChild(0).ByteRange()
		routePath := string(source[start+1 : end-1])

		start, end = node.NamedChild(1).ByteRange()
		routeMethod := string(source[start:end])

		routeProps, err := parseStruct(node.NamedChild(2), source)
		if err != nil {
			return nil, err
		}

		status := "200"
		if stat, ok := routeProps["status"]; ok {
			status = stat.(string)
		}

		routeHeaders := map[string]string{}
		routeBody := ""

		if html, ok := routeProps["html"]; ok {
			routeHeaders["Content-Type"] = "text/html"
			routeBody = html.(string)
		} else if jsonVal, ok := routeProps["json"]; ok {
			routeHeaders["Content-Type"] = "application/json"
			jsonStr, err := json.Marshal(jsonVal.(map[string]any))
			if err != nil {
				return nil, err
			}
			routeBody = string(jsonStr)
		} else if body, ok := routeProps["text"]; ok {
			routeHeaders["Content-Type"] = "text/plain"
			routeBody = body.(string)
		}

		if headers, ok := routeProps["headers"]; ok {
			routeHeaders = headers.(map[string]string)
		}

		route := route{
			Path:    routePath,
			Method:  routeMethod,
			Status:  status,
			Headers: routeHeaders,
			Body:    routeBody,
		}
		program = append(program, route)
	default:
		cursor := node.Walk()
		children := node.Children(cursor)
		for _, child := range children {
			newProgram, err := parseNode(&child, source, sourcePath, program)
			if err != nil {
				return nil, err
			}
			program = newProgram
		}
	}

	return program, nil
}

func parseStruct(node *tree_sitter.Node, source []byte) (map[string]any, error) {
	keys := []string{}
	values := []any{}

	cursor := node.Walk()
	for _, child := range node.NamedChildren(cursor) {
		switch child.GrammarName() {
		case "key":
			start, end := child.ByteRange()
			keys = append(keys, string(source[start:end]))
		case "value":
			valueType := child.NamedChild(0).GrammarName()
			switch valueType {
			case "struct":
				structVal, err := parseStruct(child.NamedChild(0), source)
				if err != nil {
					return nil, err
				}
				values = append(values, structVal)
			case "string":
				start, end := child.NamedChild(0).ByteRange()
				values = append(values, string(source[start+1:end-1]))
			default:
				start, end := child.NamedChild(0).ByteRange()
				values = append(values, string(source[start:end]))
			}
		}
	}

	if len(keys) != len(values) {
		return nil, fmt.Errorf("Mismatch between keys and values:\nkeys: %q\nvalues:%q", keys, values)
	}

	props := map[string]any{}
	for i := range keys {
		props[keys[i]] = values[i]
	}

	return props, nil
}
