package cache

import "github.com/redis/go-redis/v9"

var incrWithTTL = redis.NewScript(`
	local current = redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2])
	local ttl = redis.call('TTL', KEYS[1])

	if ttl == -1 then
		redis.call('EXPIRE', KEYS[1], 90000)
	end

	return current
`)
