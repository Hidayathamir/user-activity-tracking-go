package configkey

const (
	AESKey = "aes.key"

	AppName = "app.name"

	AuthJWTSecret        = "auth.jwt.secret"
	AuthJWTIssuer        = "auth.jwt.issuer"
	AuthJWTExpireSeconds = "auth.jwt.expire_seconds"

	DatabaseMigrations   = "database.migrations"
	DatabaseUsername     = "database.username"
	DatabasePassword     = "database.password"
	DatabaseHost         = "database.host"
	DatabasePort         = "database.port"
	DatabaseName         = "database.name"
	DatabasePoolIdle     = "database.pool.idle"
	DatabasePoolMax      = "database.pool.max"
	DatabasePoolLifetime = "database.pool.lifetime"

	LogLevel = "log.level"
)
