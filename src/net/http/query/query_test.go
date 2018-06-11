package query

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestQueryInt(t *testing.T) {
	expectedInt := 123
	req, err := http.NewRequest("GET", "/test?iamint="+strconv.Itoa(expectedInt), nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, err := q.Int("iamint", 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedInt, v)
}

func TestQueryIntByDefault(t *testing.T) {
	defaultInt := 10
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, err := q.Int("notfound", defaultInt)
	assert.NoError(t, err)
	assert.Equal(t, defaultInt, v)
}

func TestQueryIntFail(t *testing.T) {
	defaultInt := 10
	req, err := http.NewRequest("GET", "/test?iamint=abc", nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, err := q.Int("iamint", defaultInt)
	assert.Error(t, err)
	assert.Equal(t, defaultInt, v)
}

func TestInt64(t *testing.T) {
	expectedInt := int64(42949672951)
	req, err := http.NewRequest("GET", "/test?iamint="+strconv.FormatInt(expectedInt, 10), nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, err := q.Int64("iamint", 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedInt, v)
}

func TestInt64ByDefault(t *testing.T) {
	defaultInt := int64(21474836471)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, err := q.Int64("notfound", defaultInt)
	assert.NoError(t, err)
	assert.Equal(t, defaultInt, v)
}

func TestInt64Fail(t *testing.T) {
	defaultInt := int64(21474836471)
	req, err := http.NewRequest("GET", "/test?iamint=abc", nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, err := q.Int64("iamint", defaultInt)
	assert.Error(t, err)
	assert.Equal(t, defaultInt, v)
}

func TestStr(t *testing.T) {
	expectedStr := "YoYo"
	req, err := http.NewRequest("GET", "/test?Hey="+expectedStr, nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, ok := q.Str("Hey")
	assert.True(t, ok)
	assert.Equal(t, expectedStr, v)
}

func TestStrByDefault(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	q := New(req.URL.Query())

	v, ok := q.Str("notfound")
	assert.False(t, ok)
	assert.Equal(t, "", v)
}
