package cache

var _ Cache = &CacheMwLogger{}

type CacheMwLogger struct {
	Next Cache
}

func NewCacheMwLogger(next Cache) *CacheMwLogger {
	return &CacheMwLogger{
		Next: next,
	}
}
