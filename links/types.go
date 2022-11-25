package deploy

type Link struct {
	Target   string
	LinkPath string
}

type Logger interface {
	Success(message string)
	Fail(message string)
	Log(message string)
}
