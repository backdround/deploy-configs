package templates

// Template represents template to expand by this package
type Template struct {
	Name       string
	InputPath  string
	OutputPath string
	Data       interface{}
}

type Logger interface {
	Success(message string)
	Fail(message string)
	Log(message string)
}
