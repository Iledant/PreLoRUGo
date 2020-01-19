package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCoproEvent is the entry point for testing all renew projet requests
func testCoproEvent(t *testing.T, c *TestContext) {
	t.Run("CoproEvent", func(t *testing.T) {
		ID := testCreateCoproEvent(t, c)
		if ID == 0 {
			t.Error("Impossible de créer l'événement Copro")
			t.FailNow()
			return
		}
		testUpdateCoproEvent(t, c, ID)
		testGetCoproEvent(t, c, ID)
		testGetCoproEvents(t, c)
		testDeleteCoproEvent(t, c, ID)
	})
}

// testCreateCoproEvent checks if route is admin protected and created CoproEVentType
// is properly filled
func testCreateCoproEvent(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création d'événement Copro, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent: []byte(`{"CoproID":0,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création d'événement Copro : CoproID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : CoproID empty
		{
			Sent: []byte(`{"CoproID":` + strconv.FormatInt(c.CoproID, 10) +
				`,"CoproEventTypeID":0,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création d'événement Copro : CoproEventTypeID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : CoproEventTypeID empty
		{
			Sent: []byte(`{"CoproID":` + strconv.FormatInt(c.CoproID, 10) +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:  c.Config.Users.CoproUser.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"CoproEvent":{"ID":1,"CoproID":` +
				strconv.FormatInt(c.CoproID, 10) + `,"CoproEventTypeID":` +
				strconv.FormatInt(c.CoproEventTypeID, 10) + `,"Name":` +
				`"Comité d'engagement","Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`},
			StatusCode: http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro_event").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateCoproEvent", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateCoproEvent checks if route is admin protected and CoproEvent
// is properly filled
func testUpdateCoproEvent(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'événement Copro, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"CoproID":0` +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'événement Copro : CoproID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : CoproID null
		{
			Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"CoproID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"CoproEventTypeID":0,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'événement Copro : CoproEventTypeID vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : CoproID null
		{
			Sent: []byte(`{"ID":0,"CoproID":` +
				strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'événement Copro, requête : Événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{
			Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"CoproID":` +
				strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token: c.Config.Users.CoproUser.Token,
			RespContains: []string{`{"ID":` + strconv.Itoa(ID) + `,"CoproID":` +
				strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Name":"Comité d'engagement","Date":"2016-04-13T00:00:00Z","Comment":null}`},
			StatusCode: http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/copro_event").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateCoproEvent") {
		t.Error(r)
	}
}

// testGetCoproEvent checks if route is user protected and CoproEvent
// is properly filled
func testGetCoproEvent(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : no token
		{ID: 0,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Récupération d'événement Copro, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{ID: ID,
			Token: c.Config.Users.User.Token,
			RespContains: []string{`{"ID":` + strconv.Itoa(ID) + `,"CoproID":` +
				strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Name":"Comité d'engagement","Date":"2016-04-13T00:00:00Z","Comment":null}`},
			StatusCode: http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro_event/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetCoproEvent") {
		t.Error(r)
	}
}

// testGetCoproEvents checks route is protected and all CoproEvent are correctly
// sent back
func testGetCoproEvents(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Sent:  []byte(`fake`),
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"CoproEvent"`, `,"CoproID":` +
				strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"CoproEventTypeID":` + strconv.FormatInt(c.CoproEventTypeID, 10) +
				`,"Name":"Comité d'engagement","Date":"2016-04-13T00:00:00Z","Comment":null}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro_events").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetCoproEvents") {
		t.Error(r)
	}
}

// testDeleteCoproEvent checks that route is renew project protected and
// delete request sends ok back
func testDeleteCoproEvent(t *testing.T, c *TestContext, ID int) {
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
			RespContains: []string{`Suppression d'événement Copro, requête : Événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.CoproUser.Token,
			ID:           ID,
			RespContains: []string{`Événement Copro supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/copro_event/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteCoproEvent") {
		t.Error(r)
	}
}
