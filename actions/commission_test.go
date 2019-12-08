package actions

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/iris-contrib/httpexpect"
)

// testCommission is the entry point for testing all renew projet requests
func testCommission(t *testing.T, c *TestContext) {
	t.Run("Commission", func(t *testing.T) {
		ID := testCreateCommission(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la commission")
			t.FailNow()
			return
		}
		testUpdateCommission(t, c, ID)
		testGetCommission(t, c, ID)
		testGetCommissions(t, c)
		testDeleteCommission(t, c, ID)

		// Create a commission for other tests and store it's ID in context
		com := models.Commission{Name: "Commission test",
			Date: models.NullTime{Valid: true,
				Time: time.Date(2018, time.March, 1, 0, 0, 0, 0, time.UTC)}}
		if err := com.Create(c.DB); err != nil {
			t.Error("Impossible de créer la commission test : " + err.Error())
			t.FailNow()
			return
		}
		c.CommissionID = com.ID
	})
}

// testCreateCommission checks if route is admin protected and created budget action
// is properly filled
func testCreateCommission(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de commission, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Commission":{}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de commission : Champ name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"Commission":{"Name":"Essai","Date":"2019-03-01T00:00:00Z"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `{"ID"`,
			RespContains: []string{`"Commission":{"ID":1,"Name":"Essai","Date":"2019-03-01T00:00:00Z"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/commission").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateCommission", &ID)
	return ID
}

// testUpdateCommission checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateCommission(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commission, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Commission":{}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commission : Champ name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"Commission":{"ID":0,"Name":"Essai2","Date":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commission, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{
			Sent:         []byte(`{"Commission":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","Date":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Commission":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","Date":null}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/commission").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateCommission")
}

// testGetCommission checks if route is user protected and Commission correctly sent back
func testGetCommission(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de commission, requête :`},
			ID:           0}, // 1 : bad ID
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`{"Commission":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","Date":null}}`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commission/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCommission")
}

// testGetCommissions checks if route is user protected and Commissions correctly sent back
func testGetCommissions(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`{"Commission":[{"ID":1,"Name":"Essai2","Date":null}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commissions").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCommissions")
}

// testDeleteCommission checks if route is user protected and commissions correctly sent back
func testDeleteCommission(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user token
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de commission, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Commission supprimée`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/commission/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteCommission")
}
