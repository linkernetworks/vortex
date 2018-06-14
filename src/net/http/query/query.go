package query

import (
	"net/url"
	"strconv"
)

type QueryUrl struct {
	Url    url.Values    `json:"url"`
}

func New(url url.Values) *QueryUrl {
	return &QueryUrl{
		Url: url,
	}
}

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

func (query *QueryUrl) Str(key string) (string, bool) {
	values := query.Url[key]

	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}
