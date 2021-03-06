package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPEventType is the entry point for testing all renew projet requests
func testRPEventType(t *testing.T, c *TestContext) {
	t.Run("RPEventType", func(t *testing.T) {
		ID := testCreateRPEventType(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le type d'événement RP")
			t.FailNow()
			return
		}
		testUpdateRPEventType(t, c, ID)
		testGetRPEventType(t, c, ID)
		testGetRPEventTypes(t, c)
		testDeleteRPEventType(t, c, ID)
		fetchRPEventTypeID(t, c)
	})
}

// testCreateRPEventType checks if route is admin protected and created RPEVentType
// is properly filled
func testCreateRPEventType(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"Name":"Comité"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de type d'événement RP, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Name":""}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de type d'événement RP : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"Name":"Comité"}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			IDName:       `{"ID"`,
			RespContains: []string{`"RPEventType":{"ID":1,"Name":"Comité"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/rp_event_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateRPEventType", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateRPEventType checks if route is admin protected and RPEventType
// is properly filled
func testUpdateRPEventType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"Name":"Comité d'engagement"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de type d'événement RP, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Name":""}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de type d'événement RP : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"ID":0,"Name":"Comité d'engagement"}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de type d'événement RP, requête : Type d'événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"Name":"Comité d'engagement"}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`"RPEventType":{"ID":` + strconv.Itoa(ID) + `,"Name":"Comité d'engagement"}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/rp_event_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateRPEventType") {
		t.Error(r)
	}
}

// testGetRPEventType checks if route is user protected and RPEventType
// is properly filled
func testGetRPEventType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : no token
		{ID: 0,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Récupération de type d'événement RP, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{ID: ID,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"RPEventType":{"ID":` + strconv.Itoa(ID) + `,"Name":"Comité d'engagement"}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_event_type/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetRPEventType") {
		t.Error(r)
	}
}

// testGetRPEventTypes checks route is protected and all RPEventType are correctly
// sent back
func testGetRPEventTypes(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Sent:          []byte(`fake`),
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"RPEventType"`, `"Name":"Comité d'engagement"`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : bad request
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_event_types").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetRPEventTypes") {
		t.Error(r)
	}
}

// testDeleteRPEventType checks that route is renew project protected and
// delete request sends ok back
func testDeleteRPEventType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.User.Token,
			ID:           0,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 1 : user unauthorized
		{
			Token:        c.Config.Users.RenewProjectUser.Token,
			ID:           0,
			RespContains: []string{`Suppression de type d'événement RP, requête : Type d'événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.RenewProjectUser.Token,
			ID:           ID,
			RespContains: []string{`Type d'événement RP supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/rp_event_type/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteRPEventType") {
		t.Error(r)
	}
}

// fetchRPEventTypeID create an RPEventType and fetches its ID to store in the
// TestContext variable for further use
func fetchRPEventTypeID(t *testing.T, c *TestContext) {
	resp := c.E.POST("/api/rp_event_type").
		WithBytes([]byte(`{"Name":"Comité d'engagement"}`)).
		WithHeader("Authorization", "Bearer "+c.Config.Users.Admin.Token).Expect()
	body := string(resp.Content)
	status := resp.Raw().StatusCode
	if status != http.StatusCreated {
		t.Error("Impossible de créer le type d'événement pérenne")
		t.FailNow()
		return
	}
	index := strings.Index(body, `{"ID"`)
	fmt.Sscanf(body[index:], `{"ID":%d`, &c.RPEventTypeID)
	if c.RPEventTypeID == 0 {
		t.Error("Impossible de récupérer l'ID de type d'événement pérenne")
		t.FailNow()
		return
	}
}
