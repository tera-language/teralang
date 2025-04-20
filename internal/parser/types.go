package parser

type Route struct {
	Path    string
	Method  string
	Status  int
	Headers map[string]any
	Body    string
}
