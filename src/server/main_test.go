package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/config"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

// Create a test user
func init() {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)

	// init restful container
	wc := restful.NewContainer()
	service := newUserService(sp)
	wc.Add(service)

	user := entity.User{
		ID: bson.NewObjectId(),
		LoginCredential: entity.LoginCredential{
			Username: "test@linkernetworks.com",
			Password: "test",
		},
		DisplayName: "John Doe",
		Role:        "root",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0000000000",
	}

	bodyBytes, _ := json.MarshalIndent(user, "", "  ")

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, _ := http.NewRequest(
		"POST",
		"http://localhost:7890/v1/users",
		bodyReader,
	)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	wc.Dispatch(httpWriter, httpRequest)
}

func TestMain(m *testing.M) {
	v := m.Run()
	if v != 0 {
		dropTestDatabase()
	}
	os.Exit(v)
}

// Drop test database
func dropTestDatabase() bool {
	cf := config.MustRead("../../config/testing.json")
	sp := serviceprovider.NewForTesting(cf)
	session := sp.Mongo.NewSession()
	if err := session.DB("vortex_test").DropDatabase(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to drop test database.\n")
		return true
	}
	return false
}
