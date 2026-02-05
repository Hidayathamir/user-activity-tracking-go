package cachekey

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/timekit"
)

const (
	TopClientRequestCount24H           = "top_client_request_count:24h"
	TopClientRequestCount24HSourceKeys = "top_client_request_count:24h_source_keys"
)

// BuildClientRequestCountKey -.
func BuildClientRequestCountKey(apiKey string, datetime time.Time) (key string) {
	return fmt.Sprintf("client_request_count:%s:%s", apiKey, datetime.Format(time.DateOnly))
}

func ReverseClientRequestCountKey(key string) (apiKey string, datetime time.Time, err error) {
	parts := strings.Split(key, ":")

	if len(parts) != 3 {
		err := fmt.Errorf("invalid format key, key = %q", key)
		return "", time.Time{}, errkit.AddFuncName(err)
	}

	apiKey = parts[1]

	datetime, err = time.Parse(time.DateOnly, parts[2])
	if err != nil {
		return "", time.Time{}, errkit.AddFuncName(err)
	}

	return apiKey, datetime, nil
}

// BuildTopClientRequestCountHourlyKey -.
// make sure argument datetime is truncated to hour.
func BuildTopClientRequestCountHourlyKey(datetime time.Time) (key string) {
	return fmt.Sprintf("top_client_request_count:hourly:%s", datetime.Format(timekit.TimeFormatCustomRFC3339))
}

func ReverseBuildTopClientRequestCountHourlyKey(key string) (datetime time.Time, err error) {
	parts := strings.Split(key, ":")

	if len(parts) != 3 {
		err := fmt.Errorf("invalid format key, key = %q", key)
		return time.Time{}, errkit.AddFuncName(err)
	}

	datetime, err = time.Parse(timekit.TimeFormatCustomRFC3339, parts[2])
	if err != nil {
		return time.Time{}, errkit.AddFuncName(err)
	}

	return datetime, nil
}
