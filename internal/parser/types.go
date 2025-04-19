package parser

type Route struct {
	Path    string
	Method  string
	Status  string
	Headers map[string]string
	Body    string
}
