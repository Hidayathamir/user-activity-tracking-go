package counter

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
)

type KeyRequestCounter string

func (key *KeyRequestCounter) Set(apiKey string, timestamp time.Time) {
	k := fmt.Sprintf("%s||%s", apiKey, timestamp.Format(time.RFC3339))
	*key = KeyRequestCounter(k)
}

func (key *KeyRequestCounter) Parse() (apiKey string, timestamp time.Time, err error) {
	parts := strings.Split(string(*key), "||")

	if len(parts) != 2 {
		err := fmt.Errorf("invalid format key, key = %q", *key)
		return "", time.Time{}, errkit.AddFuncName(err)
	}

	apiKey = parts[0]

	timestamp, err = time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return "", time.Time{}, errkit.AddFuncName(err)
	}

	return apiKey, timestamp, nil
}

// RequestCounter is counter for insert to db.
// the key is "{apikey}||{timestamp}"
// the value is the count that in db need to increment
type RequestCounter map[KeyRequestCounter]int
