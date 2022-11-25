package deploy

type Link struct {
	TargetPath string
	LinkPath   string
}

type Logger interface {
	Success(message string)
	Fail(message string)
	Log(message string)
}
