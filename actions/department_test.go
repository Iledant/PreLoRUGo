package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testDepartment is the entry point for testing all renew projet requests
func testDepartment(t *testing.T, c *TestContext) {
	t.Run("Department", func(t *testing.T) {
		ID := testCreateDepartment(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le département")
			t.FailNow()
			return
		}
		testUpdateDepartment(t, c, ID)
		testGetDepartment(t, c, ID)
		testGetDepartments(t, c)
		testDeleteDepartment(t, c, ID)
	})
}

// testCreateDepartment checks if route is admin protected and created budget action
// is properly filled
func testCreateDepartment(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de département, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Department":{"Code":0,"Name":"Essai"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de département : Champ code incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code empty
		{
			Sent:         []byte(`{"Department":{"Code":78,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de département : Champ name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"Department":{"Code":78,"Name":"Essai"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `{"ID"`,
			RespContains: []string{`"Department":{"ID":1,"Code":78,"Name":"Essai"`},
			StatusCode:   http.StatusCreated}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/department").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateDepartment", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateDepartment checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateDepartment(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de département, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Department":{"ID":` + strconv.Itoa(ID) + `,"Code":0,"Name":"Essai2"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de département : Champ code incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code empty
		{
			Sent:         []byte(`{"Department":{"ID":` + strconv.Itoa(ID) + `,"Code":77,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de département : Champ name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"Department":{"ID":0,"Code":77,"Name":"Essai2"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de département, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{
			Sent:         []byte(`{"Department":{"ID":` + strconv.Itoa(ID) + `,"Code":77,"Name":"Essai2"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Department":{"ID":` + strconv.Itoa(ID) + `,"Code":77,"Name":"Essai2"}`},
			StatusCode:   http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/department").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateDepartment") {
		t.Error(r)
	}
}

// testGetDepartment checks if route is user protected and Department correctly sent back
func testGetDepartment(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de département, requête :`},
			ID:           0}, // 1 : bad ID
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`{"Department":{"ID":` + strconv.Itoa(ID) + `,"Code":77,"Name":"Essai2"}}`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/department/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetDepartment") {
		t.Error(r)
	}
}

// testGetDepartments checks if route is user protected and Departments correctly sent back
func testGetDepartments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`{"Department":[{"ID":1,"Code":77,"Name":"Essai2"}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/departments").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetDepartments") {
		t.Error(r)
	}
}

// testDeleteDepartment checks if route is user protected and departments correctly sent back
func testDeleteDepartment(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user token
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de département, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Département supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/department/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteDepartment") {
		t.Error(r)
	}
}
