package links

// Link is a stracture that represents symbolic link
// to create by this package.
type Link struct {
	Name       string
	TargetPath string
	LinkPath   string
}

type Logger interface {
	Success(message string)
	Fail(message string)
	Log(message string)
}
