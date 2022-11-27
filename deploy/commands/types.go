package commands

// Command represents command that creates OutputPath from
// InputPath by this package
type Command struct {
	Name       string
	InputPath  string
	OutputPath string
	Command    string
}

type Logger interface {
	Success(message string)
	Fail(message string)
	Log(message string)
}
