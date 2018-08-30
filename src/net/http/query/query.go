package query

import (
	"net/url"
	"strconv"
	"time"
)

// QueryUrl is the structure
type QueryUrl struct {
	Url url.Values `json:"url"`
}

// New will new a QueryUrl
func New(url url.Values) *QueryUrl {
	return &QueryUrl{
		Url: url,
	}
}

// Int is a function for int
func (query *QueryUrl) Int(key string, defaultValue int) (int, error) {
	values := query.Url[key]

	if len(values) == 0 {
		return defaultValue, nil
	}

	val, err := strconv.Atoi(values[0])
	if err != nil {
		return defaultValue, err
	}

	return int(val), nil
}

// Int is a function for int
func (query *QueryUrl) TimeDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	values := query.Url[key]

	if len(values) == 0 {
		return defaultValue, nil
	}

	val, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return time.Duration(val), nil
}

// Int64 is a function for int64
func (query *QueryUrl) Int64(key string, defaultValue int64) (int64, error) {
	values := query.Url[key]

	if len(values) == 0 {
		return defaultValue, nil
	}

	val, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return int64(val), nil
}

// Str is a function
func (query *QueryUrl) Str(key string) (string, bool) {
	values := query.Url[key]

	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}
