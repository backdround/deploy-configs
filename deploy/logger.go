package deploy

type Logger interface {
	Success(message string)
	Fail(message string)
	Log(message string)
}
