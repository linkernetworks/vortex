package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bitbucket.org/linkernetworks/vortex/src/entity"
	"bitbucket.org/linkernetworks/vortex/src/serviceprovider"
	restful "github.com/emicklei/go-restful"
	"github.com/influxdata/influxdb/pkg/testing/assert"
	"github.com/linkernetworks/config"
)

func TestCreateNetworkHandler(t *testing.T) {
	cf := config.MustRead("../config/testing.json")
	sp := serviceprovider.New(cf)

	network := entity.network{
		DisplayName: "OVS Bridge",
		BridgeName:  "obsbr1",
		BridgeType:  "ovs",
		Node:        "node1",
		Interface:   "eth3",
		Ports:       {2043, 2143, 2243},
		MTU:         1500,
	}

	session := sp.Mongo.NewSession()
	defer session.Close()

	bodyBytes, err := json.MarshalIndent(network, "", "  ")
	assert.NoError(t, err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://here.com/v1/networks/", bodyReader)
	assert.NoError(t, err)

	httpRequest.Header.Add("Content-Type", "application/json")

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	wc.Add(us.NewLoginService(sp))
	wc.Dispatch(httpWriter, httpRequest)

	assertResponseCode(t, 200, httpWriter)
}

// func TestEmptyFirstName(t *testing.T) {
// 	cf := config.MustRead("../config/testing.json")
// 	sp := serviceprovider.New(cf)

// 	password, err := pwdutil.EncryptPasswordLegacy("testtest")
// 	assert.NoError(t, err)

// 	user := oauth.User{
// 		ID:        bson.NewObjectId(),
// 		Email:     "test@linkernetworks.com",
// 		Password:  password,
// 		FirstName: "",
// 		LastName:  "Lin",
// 		Roles:     []string{"admin"},
// 		Verified:  true,
// 	}
// 	session := sp.Mongo.NewSession()
// 	defer session.Remove(oauth.UserCollectionName, "_id", user.ID)

// 	bodyBytes, err := json.MarshalIndent(user, "", "  ")
// 	assert.NoError(t, err)

// 	bodyReader := strings.NewReader(string(bodyBytes))
// 	httpRequest, err := http.NewRequest("POST", "http://here.com/v1/signup", bodyReader)
// 	httpRequest.Header.Add("Content-Type", "application/json")
// 	assert.NoError(t, err)

// 	httpWriter := httptest.NewRecorder()
// 	wc := restful.NewContainer()
// 	wc.Add(us.NewLoginService(sp))
// 	wc.Dispatch(httpWriter, httpRequest)
// 	assertResponseCode(t, 422, httpWriter)
// }

// func TestEmptyPassword(t *testing.T) {
// 	cf := config.MustRead("../config/testing.json")
// 	sp := serviceprovider.New(cf)

// 	password, err := pwdutil.EncryptPasswordLegacy("")
// 	assert.NoError(t, err)

// 	user := oauth.User{
// 		ID:        bson.NewObjectId(),
// 		Email:     "test@linkernetworks.com",
// 		Password:  password,
// 		FirstName: "",
// 		LastName:  "Lin",
// 		Roles:     []string{"admin"},
// 		Verified:  true,
// 	}
// 	session := sp.Mongo.NewSession()
// 	defer session.Remove(oauth.UserCollectionName, "_id", user.ID)

// 	bodyBytes, err := json.MarshalIndent(user, "", "  ")
// 	assert.NoError(t, err)

// 	bodyReader := strings.NewReader(string(bodyBytes))
// 	httpRequest, err := http.NewRequest("POST", "http://here.com/v1/signup", bodyReader)
// 	httpRequest.Header.Add("Content-Type", "application/json")
// 	assert.NoError(t, err)

// 	httpWriter := httptest.NewRecorder()
// 	wc := restful.NewContainer()
// 	wc.Add(us.NewLoginService(sp))
// 	wc.Dispatch(httpWriter, httpRequest)
// 	assertResponseCode(t, 422, httpWriter)
// }

// func TestInvalidEmail(t *testing.T) {
// 	cf := config.MustRead("../config/testing.json")
// 	sp := serviceprovider.New(cf)

// 	password, err := pwdutil.EncryptPasswordLegacy("")
// 	assert.NoError(t, err)

// 	user := oauth.User{
// 		ID:        bson.NewObjectId(),
// 		Email:     "",
// 		Password:  password,
// 		FirstName: "Tester",
// 		LastName:  "Lin",
// 		Roles:     []string{"admin"},
// 		Verified:  true,
// 	}
// 	session := sp.Mongo.NewSession()
// 	defer session.Remove(oauth.UserCollectionName, "_id", user.ID)

// 	bodyBytes, err := json.MarshalIndent(user, "", "  ")
// 	assert.NoError(t, err)

// 	bodyReader := strings.NewReader(string(bodyBytes))
// 	httpRequest, err := http.NewRequest("POST", "http://here.com/v1/signup", bodyReader)
// 	httpRequest.Header.Add("Content-Type", "application/json")
// 	assert.NoError(t, err)

// 	httpWriter := httptest.NewRecorder()
// 	wc := restful.NewContainer()
// 	wc.Add(us.NewLoginService(sp))
// 	wc.Dispatch(httpWriter, httpRequest)
// 	assertResponseCode(t, 422, httpWriter)
// }
