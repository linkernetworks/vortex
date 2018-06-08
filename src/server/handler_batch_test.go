package server

import (
	"testing"

	"net/http/httptest"
)

func assertResponseCode(t *testing.T, expectedCode int, resp *httptest.ResponseRecorder) {
	t.Helper()
	t.Logf("code:%d", resp.Code)
	if expectedCode != resp.Code {
		t.Errorf("status code %d expected.", expectedCode)
		t.Logf("Response:\n%s", resp.Body.String())
	}
}
