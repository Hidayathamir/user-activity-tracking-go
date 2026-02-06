package column

import (
	"fmt"
)

type Column string

func (c Column) Str() string {
	return string(c)
}

func (c Column) SumAs(as string) string {
	return fmt.Sprintf("SUM(%s) as %s", string(c), as)
}

func (c Column) Eq(value any) (string, any) {
	return string(c) + " = ?", value
}

func (c Column) GTE(value any) (string, any) {
	return string(c) + " >= ?", value
}

func (c Column) LT(value any) (string, any) {
	return string(c) + " < ?", value
}

func (c Column) Plus(value any) (string, any) {
	return string(c) + " + ?", value
}

func (c Column) Desc() string {
	return string(c) + " DESC"
}

const (
	ID        Column = "id"
	Name      Column = "name"
	Email     Column = "email"
	Password  Column = "password"
	APIKey    Column = "api_key"
	IP        Column = "ip"
	Endpoint  Column = "endpoint"
	Timestamp Column = "timestamp"
	Datetime  Column = "datetime"
	Count     Column = "count"
)
