package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPEvent is the entry point for testing all renew projet requests
func testRPEvent(t *testing.T, c *TestContext) {
	t.Run("RPEvent", func(t *testing.T) {
		ID := testCreateRPEvent(t, c)
		if ID == 0 {
			t.Error("Impossible de créer l'événement RP")
			t.FailNow()
			return
		}
		testUpdateRPEvent(t, c, ID)
		testGetRPEvent(t, c, ID)
		testGetRPEvents(t, c)
		testDeleteRPEvent(t, c, ID)
	})
}

// testCreateRPEvent checks if route is admin protected and created RPEVentType
// is properly filled
func testCreateRPEvent(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"RenewProjectID":` + strconv.FormatInt(c.RenewProjectID, 10) +
			`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création d'événement RP, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"RenewProjectID":0,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création d'événement RP : Champ RenewProjectID ou RPEventTypeID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : RenewProjectID empty
		{Sent: []byte(`{"RenewProjectID":` + strconv.FormatInt(c.RenewProjectID, 10) +
			`,"RPEventTypeID":0,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création d'événement RP : Champ RenewProjectID ou RPEventTypeID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : RPEventTypeID empty
		{Sent: []byte(`{"RenewProjectID":` + strconv.FormatInt(c.RenewProjectID, 10) +
			`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`),
			Token:  c.Config.Users.RenewProjectUser.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"RPEvent":{"ID":1,"RenewProjectID":` +
				strconv.FormatInt(c.RenewProjectID, 10) + `,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"Date":"2015-04-13T00:00:00Z","Comment":"Commentaire"}`},
			StatusCode: http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/rp_event").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateRPEvent", &ID)
	return ID
}

// testUpdateRPEvent checks if route is admin protected and RPEvent
// is properly filled
func testUpdateRPEvent(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"RenewProjectID":` +
			strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification d'événement RP, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"RenewProjectID":0` +
			`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification d'événement RP : Champ RenewProjectID ou RPEventTypeID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : RenewProjectID null
		{Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"RenewProjectID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"RPEventTypeID":0,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification d'événement RP : Champ RenewProjectID ou RPEventTypeID vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : RenewProjectID null
		{Sent: []byte(`{"ID":0,"RenewProjectID":` +
			strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification d'événement RP, requête : Événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"RenewProjectID":` +
			strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
			`,"Date":"2016-04-13T00:00:00Z","Comment":null}`),
			Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`{"ID":` + strconv.Itoa(ID) + `,"RenewProjectID":` +
				strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"Date":"2016-04-13T00:00:00Z","Comment":null}`},
			StatusCode: http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/rp_event").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateRPEvent")
}

// testGetRPEvent checks if route is user protected and RPEvent
// is properly filled
func testGetRPEvent(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Token:        "",
			RespContains: []string{`Token absent`},
			StatusCode:   http.StatusInternalServerError}, // 0 : no token
		{ID: 0,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Récupération d'événement RP, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{ID: ID,
			Token: c.Config.Users.User.Token,
			RespContains: []string{`{"ID":` + strconv.Itoa(ID) + `,"RenewProjectID":` +
				strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"Date":"2016-04-13T00:00:00Z","Comment":null}`},
			StatusCode: http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_event/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPEvent")
}

// testGetRPEvents checks route is protected and all RPEvent are correctly
// sent back
func testGetRPEvents(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "fake",
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"RPEvent"`, `,"RenewProjectID":` +
				strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"RPEventTypeID":` + strconv.FormatInt(c.RPEventTypeID, 10) +
				`,"Date":"2016-04-13T00:00:00Z","Comment":null}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_events").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPEvents")
}

// testDeleteRPEvent checks that route is renew project protected and
// delete request sends ok back
func testDeleteRPEvent(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "fake",
			ID:           0,
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : bad token
		{Token: c.Config.Users.User.Token,
			ID:           0,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 1 : user unauthorized
		{Token: c.Config.Users.RenewProjectUser.Token,
			ID:           0,
			RespContains: []string{`Suppression d'événement RP, requête : Événement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{Token: c.Config.Users.RenewProjectUser.Token,
			ID:           ID,
			RespContains: []string{`Événement RP supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/rp_event/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteRPEvent")
}
