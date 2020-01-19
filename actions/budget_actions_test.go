package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBudgetAction is the entry point for testing all budget action related routes
func testBudgetAction(t *testing.T, c *TestContext) {
	t.Run("BudgetAction", func(t *testing.T) {
		ID := testCreateBudgetAction(t, c)
		if ID == 0 {
			t.Error("Impossible de créer l'action budgétaire")
			t.FailNow()
			return
		}
		testUpdateBudgetAction(t, c, ID)
		testGetBudgetActions(t, c)
		testDeleteBudgetAction(t, c, ID)
	})
}

// testCreateBudgetAction checks if route is admin protected and created budget action
// is properly filled
func testCreateBudgetAction(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'action budgétaire, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : code empty
		{
			Sent:         []byte(`{"BudgetAction":{"Code":0,"Name":"Action"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code empty
		{
			Sent:         []byte(`{"BudgetAction":{"Code":1234567890,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"BudgetAction":{"Code":1234567890,"Name":"Action"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `{"ID"`,
			RespContains: []string{`"BudgetAction"`, `"Code":1234567890,"Name":"Action"`},
			StatusCode:   http.StatusCreated}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/budget_action").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateBudgetAction", &ID) {
		t.Error(r)
	}
	return ID
}

// testGetBudgetActions checks route is protected and datas sent back are well formed
func testGetBudgetActions(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token null
		{
			Token:         c.Config.Users.Admin.Token,
			RespContains:  []string{`BudgetAction`, `Name`, `Code`, `BudgetSector`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/budget_actions").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBudgetActions") {
		t.Error(r)
	}
}

// testUpdateBudgetAction checks if route is admin protected and budget action
// is properly modified
func testUpdateBudgetAction(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : code empty
		{
			Sent:         []byte(`{"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":0,"Name":"Action"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code zero
		{
			Sent:         []byte(`{"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":1234567890,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"BudgetAction":{"ID":0,"Code":1234567890,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : name empty
		{
			Sent:         []byte(`{"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":23456789,"Name":"Action modifiée"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":23456789,"Name":"Action modifiée","SectorID":null}}`},
			StatusCode:   http.StatusOK}, // 5 : name empty
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/budget_action").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateBudgetAction") {
		t.Error(r)
	}
}

// testDeleteBudgetAction check route is admin protected and delete requests returns ok
func testDeleteBudgetAction(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression d'action budgétaire, requête : Action budgétaire introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Action budgétaire supprimée`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/budget_action/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteBudgetAction") {
		t.Error(r)
	}
}
