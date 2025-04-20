package parser

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path"
	"slices"
	"strconv"

	"github.com/tera-language/teralang/internal/logger"
	tree_sitter_teralang "github.com/tera-language/tree-sitter-teralang/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

var parsedPaths []string = []string{}

func Parse(filepath string) ([]Route, error) {
	logger.Infoln("Parsing file:", filepath)

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

	program := []Route{}
	program, err = parseNode(rootNode, text, filepath, program)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func parseNode(node *tree_sitter.Node, source []byte, sourcePath string, program []Route) ([]Route, error) {
	switch node.GrammarName() {
	case "import":
		start, end := node.NamedChild(0).ByteRange()
		importPath := string(source[start+1 : end-1]) // skip "" characters
		relativePath := path.Join(path.Dir(sourcePath), importPath)
		if slices.Contains(parsedPaths, relativePath) {
			logger.Warningln("Already parsed:", relativePath)
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
		keys := maps.Keys(routeProps)
		for key := range keys {
			if !slices.Contains([]string{"status", "headers", "json", "html", "text"}, key) {
				return nil, fmt.Errorf("Invalid key in route definition: \"%s\"", key)
			}
		}

		status := 200
		if stat, ok := routeProps["status"]; ok {
			status, ok = stat.(int)
			if !ok {
				return nil, fmt.Errorf("Error on line %d: \"%v\" should be an integer", node.StartPosition().Row+1, stat)
			}
		}

		routeHeaders := map[string]any{}
		routeBody := ""

		if html, ok := routeProps["html"]; ok {
			routeHeaders["Content-Type"] = "text/html"
			routeBody = html.(string)
		} else if jsonVal, ok := routeProps["json"]; ok {
			routeHeaders["Content-Type"] = "application/json"
			jsonMap := jsonVal.(map[string]any)
			jsonStr, err := json.Marshal(jsonMap)
			if err != nil {
				return nil, fmt.Errorf("Error on line %d: %s", node.StartPosition().Row+1, err)
			}
			routeBody = string(jsonStr)
		} else if body, ok := routeProps["text"]; ok {
			routeHeaders["Content-Type"] = "text/plain"
			routeBody, ok = body.(string)
		}

		if headers, ok := routeProps["headers"]; ok {
			routeHeaders = headers.(map[string]any)
		}

		route := Route{
			Path:    routePath,
			Method:  routeMethod,
			Status:  int(status),
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
	startPos := []uint{}

	cursor := node.Walk()
	for _, child := range node.NamedChildren(cursor) {
		switch child.GrammarName() {
		case "key":
			if child.NamedChildCount() == 1 {
				child = *child.NamedChild(0)
				start, end := child.ByteRange()
				keys = append(keys, string(source[start+1:end-1]))
			} else {
				start, end := child.ByteRange()
				keys = append(keys, string(source[start:end]))
			}
			startPos = append(startPos, child.StartPosition().Row+1)
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
			case "int":
				start, end := child.NamedChild(0).ByteRange()
				srcStr := string(source[start:end])
				intVal, err := strconv.ParseInt(srcStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("Error on line %d: %v", child.StartPosition().Row, err)
				}
				values = append(values, int(intVal))
			case "float":
				start, end := child.NamedChild(0).ByteRange()
				srcStr := string(source[start:end])
				floatVal, err := strconv.ParseFloat(srcStr, 64)
				if err != nil {
					return nil, fmt.Errorf("Error on line %d: %v", child.StartPosition().Row, err)
				}
				values = append(values, floatVal)
			case "bool":
				start, end := child.NamedChild(0).ByteRange()
				srcStr := string(source[start:end])
				boolVal := false
				if srcStr == "true" {
					boolVal = true
				}
				values = append(values, boolVal)
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
		switch keys[i] {
		case "headers":
			if _, ok := values[i].(map[string]any); !ok {
				return nil, fmt.Errorf("Error on line %d: \"headers\" should be a struct", startPos[i])
			}
		case "json":
			if _, ok := values[i].(map[string]any); !ok {
				return nil, fmt.Errorf("Error on line %d: \"json\" should be a struct", startPos[i])
			}
		case "text":
			if _, ok := values[i].(string); !ok {
				return nil, fmt.Errorf("Error on line %d: \"text\" should be a string", startPos[i])
			}
		}

		props[keys[i]] = values[i]
	}

	return props, nil
}
