package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBudgetSector is the entry point for testing all renew projet requests
func testBudgetSector(t *testing.T, c *TestContext) {
	t.Run("BudgetSector", func(t *testing.T) {
		ID := testCreateBudgetSector(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le secteur budgétaire")
			t.FailNow()
			return
		}
		testUpdateBudgetSector(t, c, ID)
		testGetBudgetSector(t, c, ID)
		testGetBudgetSectors(t, c)
		testDeleteBudgetSector(t, c, ID)
	})
}

// testCreateBudgetSector checks if route is admin protected and created budget action
// is properly filled
func testCreateBudgetSector(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"BudgetSector":{"Name":"Essai","FullName":"Essai"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de secteur budgétaire, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"BudgetSector":{"Name":"","FullName":"Essai"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de secteur budgétaire : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : empty name
		{Sent: []byte(`{"BudgetSector":{"Name":"Essai","FullName":"Essai"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetSector":{"ID":2,"Name":"Essai","FullName":"Essai"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/budget_sector").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateBudgetSector", &ID)
	return ID
}

// testUpdateBudgetSector checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateBudgetSector(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"BudgetSector":{"Name":"Essai2","FullName":null}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de secteur budgétaire, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"BudgetSector":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de secteur budgétaire : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{Sent: []byte(`{"BudgetSector":{"ID":0,"Name":"Essai2","FullName":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de secteur budgétaire, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{Sent: []byte(`{"BudgetSector":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","FullName":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetSector":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","FullName":null}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/budget_sector").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateBudgetSector")
}

// testGetBudgetSector checks if route is user protected and BudgetSector correctly sent back
func testGetBudgetSector(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de secteur budgétaire, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"BudgetSector":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","FullName":null}}`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/budget_sector/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetBudgetSector")
}

// testGetBudgetSectors checks if route is user protected and BudgetSectors correctly sent back
func testGetBudgetSectors(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        2,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains:  []string{`{"BudgetSector":[{"ID":1,"Name":"LO","FullName":null},{"ID":2,"Name":"Essai2","FullName":null}]}`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/budget_sectors").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetBudgetSectors")
}

// testDeleteBudgetSector checks if route is user protected and budget_sectors correctly sent back
func testDeleteBudgetSector(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de secteur budgétaire, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Logement supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/budget_sector/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteBudgetSector")
}
