package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingConvention is the entry point for testing all renew projet requests
func testHousingConvention(t *testing.T, c *TestContext) {
	t.Run("HousingConvention", func(t *testing.T) {
		ID := testCreateHousingConvention(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la convention de logement")
			t.FailNow()
			return
		}
		testUpdateHousingConvention(t, c, ID)
		testGetHousingConventions(t, c)
		testDeleteHousingConvention(t, c, ID)
	})
}

// testCreateHousingConvention checks if route is admin protected and created housing
// convention is properly filled
func testCreateHousingConvention(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de convention de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingConvention":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de convention de logement, paramètre :`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"HousingConvention":{"Name":"T3"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"HousingConvention"`, `"Name":"T3"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_convention").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingConvention", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{{
		Sent:         []byte(`{"HousingConvention":{"Name":"T5"}}`),
		Token:        c.Config.Users.Admin.Token,
		IDName:       `"ID"`,
		RespContains: []string{`"HousingConvention"`, `"Name":"T5"`},
		StatusCode:   http.StatusCreated}, // created to store in config
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingConventionForTests", &c.HousingConventionID) {
		t.Error(r)
	}
	return ID
}

// testUpdateHousingConvention checks if route is admin protected and housing convention
// is properly filled
func testUpdateHousingConvention(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de convention de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingConvention":{"ID":0,"Name":"T4"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de convention de logement, requête : Convention introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"HousingConvention":{"ID":` + strconv.Itoa(ID) + `,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de convention de logement, paramètre : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"HousingConvention":{"ID":` + strconv.Itoa(ID) + `,"Name":"T4"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"HousingConvention"`, `"Name":"T4"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/housing_convention").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateHousingConvention") {
		t.Error(r)
	}
}

// testGetHousingConventions checks route is protected and all housing conventions
// are correctly sent back
func testGetHousingConventions(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"HousingConvention"`, `"Name"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_conventions").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetHousingConventions") {
		t.Error(r)
	}
}

// testDeleteHousingConvention checks that route is renew project protected and
// delete request sends ok back
func testDeleteHousingConvention(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de convention de logement, requête : Convention introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Convention de logement supprimée`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/housing_convention/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteHousingConvention") {
		t.Error(r)
	}
}
