package repository

var _ ClientRequestCountRepository = &ClientRequestCountRepositoryMwLogger{}

type ClientRequestCountRepositoryMwLogger struct {
	Next ClientRequestCountRepository
}

func NewClientRequestCountRepositoryMwLogger(next ClientRequestCountRepository) *ClientRequestCountRepositoryMwLogger {
	return &ClientRequestCountRepositoryMwLogger{
		Next: next,
	}
}
