package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingType is the entry point for testing all housing types
func testHousingType(t *testing.T, c *TestContext) {
	t.Run("HousingType", func(t *testing.T) {
		ID := testCreateHousingType(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le type de logement")
			t.FailNow()
			return
		}
		testUpdateHousingType(t, c, ID)
		testGetHousingTypes(t, c)
		testDeleteHousingType(t, c, ID)
	})
}

// testCreateHousingType checks if route is admin protected and created housing
// type is properly filled
func testCreateHousingType(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de type de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingType":{"ShortName":"","LongName":"type de logement"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de type de logement, paramètre :`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"HousingType":{"ShortName":"TL","LongName":"type de logement"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"HousingType"`, `"ShortName":"TL"`, `"LongName":"type de logement"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingType", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{
		{
			Sent:       []byte(`{"HousingType":{"ShortName":"LF","LongName":"Logement familial"}}`),
			Token:      c.Config.Users.Admin.Token,
			IDName:     `"ID"`,
			StatusCode: http.StatusCreated}, // 0 : create housing type
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingTypeID", &c.HousingTypeID) {
		t.Error(r)
	}
	return ID
}

// testUpdateHousingType checks if route is admin protected and housing type
// is properly filled
func testUpdateHousingType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de type de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingType":{"ID":0,"ShortName":"TML","LongName":"type modifié de logement"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de type de logement, requête : Type introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"HousingType":{"ID":` + strconv.Itoa(ID) + `,"ShortName":"","LongName":"type modifié de logement"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de type de logement, paramètre : Nom court vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"HousingType":{"ID":` + strconv.Itoa(ID) + `,"ShortName":"TML","LongName":"type modifié de logement"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"HousingType"`, `"ShortName":"TML"`, `"LongName":"type modifié de logement"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/housing_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateHousingType") {
		t.Error(r)
	}
}

// testGetHousingTypes checks route is protected and all housing types
// are correctly sent back
func testGetHousingTypes(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"HousingType":[`, `"ShortName"`, `"LongName"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_types").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetHousingTypes") {
		t.Error(r)
	}
}

// testDeleteHousingType checks that route is renew project protected and
// delete request sends ok back
func testDeleteHousingType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de type de logement, requête : Type introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Type de logement supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/housing_type/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteHousingType") {
		t.Error(r)
	}
}
