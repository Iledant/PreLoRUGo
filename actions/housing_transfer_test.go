package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingTransfer is the entry point for testing all housing comments
func testHousingTransfer(t *testing.T, c *TestContext) {
	t.Run("HousingTransfer", func(t *testing.T) {
		ID := testCreateHousingTransfer(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le transfert de logement")
			t.FailNow()
			return
		}
		testUpdateHousingTransfer(t, c, ID)
		testGetHousingTransfers(t, c)
		testDeleteHousingTransfer(t, c, ID)
	})
}

// testCreateHousingTransfer checks if route is admin protected and created housing
// comment is properly filled
func testCreateHousingTransfer(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de transfert de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingTransfer":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de transfert de logement, paramètre :`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"HousingTransfer":{"Name":"type de transfert"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"HousingTransfer"`, `"Name":"type de transfert"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_transfer").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingTransfer", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{{
		Sent:         []byte(`{"HousingTransfer":{"Name":"type de transfert test"}}`),
		Token:        c.Config.Users.Admin.Token,
		IDName:       `"ID"`,
		RespContains: []string{`"HousingTransfer"`, `"Name":"type de transfert test"`},
		StatusCode:   http.StatusCreated}, // created to store in config
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingTransferForTests", &c.HousingTransferID) {
		t.Error(r)
	}
	return ID
}

// testUpdateHousingTransfer checks if route is admin protected and housing comment
// is properly filled
func testUpdateHousingTransfer(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de transfert de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingTransfer":{"ID":0,"Name":"type modifié de transfert"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de transfert de logement, requête : Transfert introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"HousingTransfer":{"ID":` + strconv.Itoa(ID) + `,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de transfert de logement, paramètre : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"HousingTransfer":{"ID":` + strconv.Itoa(ID) + `,"Name":"type modifié de transfert"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"HousingTransfer"`, `"Name":"type modifié de transfert"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/housing_transfer").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateHousingTransfer") {
		t.Error(r)
	}
}

// testGetHousingTransfers checks route is protected and all housing conventions
// are correctly sent back
func testGetHousingTransfers(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"HousingTransfer"`, `"Name"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_transfers").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetHousingTransfers") {
		t.Error(r)
	}
}

// testDeleteHousingTransfer checks that route is renew project protected and
// delete request sends ok back
func testDeleteHousingTransfer(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de transfert de logement, requête : Transfert introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Transfert de logement supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/housing_transfer/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteHousingTransfer") {
		t.Error(r)
	}
}
