package cache

import (
	"errors"
	"net"

	"github.com/redis/go-redis/v9"
)

func isRedisUnavailable(err error) bool {
	if err == nil {
		return false
	}

	// Connection refused, network unreachable, etc.
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return true
	}

	// Redis client explicitly closed
	if errors.Is(err, redis.ErrClosed) {
		return true
	}

	return false
}
