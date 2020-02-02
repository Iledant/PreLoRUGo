package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingTypology is the entry point for testing all renew projet requests
func testHousingTypology(t *testing.T, c *TestContext) {
	t.Run("HousingTypology", func(t *testing.T) {
		ID := testCreateHousingTypology(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la typologie de logement")
			t.FailNow()
			return
		}
		testUpdateHousingTypology(t, c, ID)
		testGetHousingTypologies(t, c)
		testDeleteHousingTypology(t, c, ID)
	})
}

// testCreateHousingTypology checks if route is admin protected and created housing
// typology is properly filled
func testCreateHousingTypology(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de typologie de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingTypology":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de typologie de logement, paramètre :`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"HousingTypology":{"Name":"T3"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"HousingTypology"`, `"Name":"T3"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_typology").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingTypology", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{{
		Sent:         []byte(`{"HousingTypology":{"Name":"T5"}}`),
		Token:        c.Config.Users.Admin.Token,
		IDName:       `"ID"`,
		RespContains: []string{`"HousingTypology"`, `"Name":"T5"`},
		StatusCode:   http.StatusCreated}, // created to store in config
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingTypologyForTests", &c.HousingTypologyID) {
		t.Error(r)
	}
	return ID
}

// testUpdateHousingTypology checks if route is admin protected and HousingTypology
// is properly filled
func testUpdateHousingTypology(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de typologie de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingTypology":{"ID":0,"Name":"T4"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de typologie de logement, requête : Typologie introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"HousingTypology":{"ID":` + strconv.Itoa(ID) + `,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de typologie de logement, paramètre : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"HousingTypology":{"ID":` + strconv.Itoa(ID) + `,"Name":"T4"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"HousingTypology"`, `"Name":"T4"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/housing_typology").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateHousingTypology") {
		t.Error(r)
	}
}

// testGetHousingTypologies checks route is protected and all HousingTypology are correctly
// sent back
func testGetHousingTypologies(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"HousingTypology"`, `"Name"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_typologies").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetHousingTypologies") {
		t.Error(r)
	}
}

// testDeleteHousingTypology checks that route is renew project protected and
// delete request sends ok back
func testDeleteHousingTypology(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de typologie de logement, requête : Typologie introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Typologie de logement supprimée`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/housing_typology/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteHousingTypology") {
		t.Error(r)
	}
}
