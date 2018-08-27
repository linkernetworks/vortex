package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/entity"
	response "github.com/linkernetworks/vortex/src/net/http"
)

func loginGetToken(wc *restful.Container) (string, error) {
	var resp response.ActionResponse
	userCred := entity.LoginCredential{
		Username: "test@linkernetworks.com",
		Password: "test",
	}

	bodyBytes, err := json.MarshalIndent(userCred, "", "  ")
	if err != nil {
		return "", err
	}

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest(
		"POST",
		"http://localhost:7890/v1/users/signin",
		bodyReader,
	)
	if err != nil {
		return "", err
	}

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)

	decoder := json.NewDecoder(httpWriter.Body)
	if err := decoder.Decode(&resp); err != nil {
		return "", err
	}
	return string(resp.Message), nil
}
