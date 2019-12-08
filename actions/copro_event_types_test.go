package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCoproEventType is the entry point for testing all renew projet requests
func testCoproEventType(t *testing.T, c *TestContext) {
	t.Run("CoproEventType", func(t *testing.T) {
		ID := testCreateCoproEventType(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le type d'événement Copro")
			t.FailNow()
			return
		}
		testUpdateCoproEventType(t, c, ID)
		testGetCoproEventType(t, c, ID)
		testGetCoproEventTypes(t, c)
		testDeleteCoproEventType(t, c, ID)
		fetchCoproEventTypeID(t, c)
	})
}

// testCreateCoproEventType checks if route is admin protected and created CoproEVentType
// is properly filled
func testCreateCoproEventType(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création de type d'événement Copro, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Name":""}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création de type d'événement Copro : name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"Name":"Comité"}`),
			Token:        c.Config.Users.CoproUser.Token,
			IDName:       `{"ID"`,
			RespContains: []string{`"CoproEventType":{"ID":1,"Name":"Comité"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro_event_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateCoproEventType", &ID)
	return ID
}

// testUpdateCoproEventType checks if route is admin protected and CoproEventType
// is properly filled
func testUpdateCoproEventType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de type d'événement Copro, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Name":""}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de type d'événement Copro : name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"ID":0,"Name":"Comité d'engagement"}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de type d'événement Copro, requête : Type d'événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"Name":"Comité d'engagement"}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`"CoproEventType":{"ID":` + strconv.Itoa(ID) + `,"Name":"Comité d'engagement"}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/copro_event_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateCoproEventType")
}

// testGetCoproEventType checks if route is user protected and CoproEventType
// is properly filled
func testGetCoproEventType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : no token
		{
			ID:           0,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Récupération de type d'événement Copro, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			ID:           ID,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"CoproEventType":{"ID":` + strconv.Itoa(ID) + `,"Name":"Comité d'engagement"}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro_event_type/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoproEventType")
}

// testGetCoproEventTypes checks route is protected and all CoproEventType are correctly
// sent back
func testGetCoproEventTypes(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Sent:          []byte(`fake`),
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"CoproEventType"`, `"Name":"Comité d'engagement"`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : bad request
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro_event_types").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoproEventTypes")
}

// testDeleteCoproEventType checks that route is renew project protected and
// delete request sends ok back
func testDeleteCoproEventType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.User.Token,
			ID:           0,
			RespContains: []string{`Droits sur les copropriétés requis`},
			StatusCode:   http.StatusUnauthorized}, // 1 : user unauthorized
		{
			Token:        c.Config.Users.CoproUser.Token,
			ID:           0,
			RespContains: []string{`Suppression de type d'événement Copro, requête : Type d'événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.CoproUser.Token,
			ID:           ID,
			RespContains: []string{`Type d'événement Copro supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/copro_event_type/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteCoproEventType")
}

// fetchCoproEventTypeID create an CoproEventType and fetches its ID to store in the
// TestContext variable for further use
func fetchCoproEventTypeID(t *testing.T, c *TestContext) {
	resp := c.E.POST("/api/copro_event_type").
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
	fmt.Sscanf(body[index:], `{"ID":%d`, &c.CoproEventTypeID)
	if c.CoproEventTypeID == 0 {
		t.Error("Impossible de récupérer l'ID de type d'événement pérenne")
		t.FailNow()
		return
	}
}
