package repository

var _ RequestLogRepository = &RequestLogRepositoryMwLogger{}

type RequestLogRepositoryMwLogger struct {
	Next RequestLogRepository
}

func NewRequestLogRepositoryMwLogger(next RequestLogRepository) *RequestLogRepositoryMwLogger {
	return &RequestLogRepositoryMwLogger{
		Next: next,
	}
}
