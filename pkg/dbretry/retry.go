package dbretry

import (
	"errors"
	"strings"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Default retry configuration
const (
	DefaultAttempts = 3
	DefaultDelay    = 500 * time.Millisecond
)

// Do wraps a database operation with retry logic.
// It only retries transient/recoverable errors like deadlocks or connection issues.
// Non-retryable errors (duplicate key, foreign key violations, etc.) are returned immediately.
func Do(fn func() error, opts ...retry.Option) error {
	defaultOpts := []retry.Option{
		retry.Attempts(DefaultAttempts),
		retry.Delay(DefaultDelay),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(IsRetryable),
		retry.OnRetry(func(n uint, err error) {
			x.Logger.WithError(err).WithField("attempt", n+1).Warn("query db go error, db is retrying")
		}),
	}

	// User options override defaults
	allOpts := append(defaultOpts, opts...)

	return retry.Do(fn, allOpts...)
}

// DoWithResult wraps a database operation that returns a value with retry logic.
func DoWithResult[T any](fn func() (T, error), opts ...retry.Option) (T, error) {
	defaultOpts := []retry.Option{
		retry.Attempts(DefaultAttempts),
		retry.Delay(DefaultDelay),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(IsRetryable),
		retry.OnRetry(func(n uint, err error) {
			x.Logger.WithError(err).WithField("attempt", n+1).Warn("query db go error, db is retrying")
		}),
	}

	allOpts := append(defaultOpts, opts...)

	return retry.DoWithData(fn, allOpts...)
}

// IsRetryable checks if the error is transient and worth retrying.
// Returns true for deadlocks, connection issues, and other transient errors.
// Returns false for constraint violations, record not found, and other permanent errors.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Don't retry GORM's known non-transient errors
	if errors.Is(err, gorm.ErrDuplicatedKey) ||
		errors.Is(err, gorm.ErrForeignKeyViolated) ||
		errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, gorm.ErrCheckConstraintViolated) ||
		errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidValue) {
		return false
	}

	// Check PostgreSQL error codes for retryable conditions
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40001": // serialization_failure
			return true
		case "40P01": // deadlock_detected
			return true
		case "57P01": // admin_shutdown
			return true
		case "57P02": // crash_shutdown
			return true
		case "57P03": // cannot_connect_now
			return true
		case "08000": // connection_exception
			return true
		case "08003": // connection_does_not_exist
			return true
		case "08006": // connection_failure
			return true
		}
		// Don't retry other PG errors (constraint violations, syntax, etc.)
		return false
	}

	// Check for connection-related errors by message
	errMsg := strings.ToLower(err.Error())
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"timeout",
		"no connection",
		"connection closed",
		"driver: bad connection",
		"invalid connection",
	}
	for _, pattern := range retryablePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
