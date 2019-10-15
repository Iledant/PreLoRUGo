package actions

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/iris-contrib/httpexpect"
)

// testCopro is the entry point for testing all copro routes
func testCopro(t *testing.T, c *TestContext) {
	t.Run("Copro", func(t *testing.T) {
		ID := testCreateCopro(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la copropriété")
			t.FailNow()
			return
		}
		testModifyCopro(t, c, ID)
		testGetCopros(t, c)
		testGetCoproDatas(t, c, ID)
		testDeleteCopro(t, c, ID)
		testBatchCopros(t, c)
		copro := models.Copro{Reference: "RefCoproTest",
			Name:      "Copro Test",
			Address:   models.NullString{Valid: true, String: "Adresse de test"},
			ZipCode:   models.NullInt64{Valid: true, Int64: 77001},
			LabelDate: models.NullTime{Valid: false},
			Budget:    models.NullInt64{Valid: false}}
		if err := copro.Create(c.DB); err != nil {
			t.Error("Impossible de créer la copropriété de test")
			t.FailNow()
			return
		}
		c.CoproID = copro.ID
	})
}

// testCreateCopro check if route is protected and copro correctly created
func testCreateCopro(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : no token
		{Sent: []byte(`{"Copro":{"Reference":"","Name":"","Address":"","ZipCode":0,` +
			`"LabelDate":null,"Budget":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Création de copropriété : champ reference vide"`},
			StatusCode:   http.StatusBadRequest}, // 1 : reference empty
		{Sent: []byte(`{"Copro":{"Reference":"CO001","Name":"","Address":"",` +
			`"ZipCode":0,"LabelDate":null,"Budget":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Création de copropriété : champ name vide"`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{Sent: []byte(`{"Copro":{"Reference":"CO001","Name":"Copro","Address":"adresse",` +
			`"ZipCode":93200,"LabelDate":"2016-03-01T12:00:00Z","Budget":1000000}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de copropriété, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 3 : zipcode doesn't exist
		{Sent: []byte(`{"Copro":{"Reference":"CO001","Name":"Copro","Address":"adresse",` +
			`"ZipCode":78146,"LabelDate":"2016-03-01T12:00:00Z","Budget":1000000}}`),
			Token:  c.Config.Users.Admin.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"Copro":{"ID":2,"Reference":"CO001","Name":"Copro",` +
				`"Address":"adresse","ZipCode":78146,"CityName":"CHATOU",` +
				`"LabelDate":"2016-03-01T12:00:00Z","Budget":1000000`},
			StatusCode: http.StatusCreated}, // 4 : ok
		{Sent: []byte(`{"Copro":{"Reference":"CO100","Name":"Copro cadre","Address":null,` +
			`"ZipCode":null,"LabelDate":null,"Budget":null}}`),
			Token:  c.Config.Users.Admin.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"Copro":{"ID":3,"Reference":"CO100","Name":"Copro cadre",` +
				`"Address":null,"ZipCode":null,"CityName":null,` +
				`"LabelDate":null,"Budget":null`},
			StatusCode: http.StatusCreated}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateCopro", &ID)
	return ID
}

// testModifyCopro check route is protected for admin and modifications are correctly done
func testModifyCopro(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : no token
		{Sent: []byte(`{"Copro":{"Reference":"","Name":"","Address":"","ZipCode":0,` +
			`"LabelDate":null,"Budget":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Modification de copropriété : champ reference vide"`},
			StatusCode:   http.StatusBadRequest}, // 1 : reference empty
		{Sent: []byte(`{"Copro":{"Reference":"CO001","Name":"","Address":"","ZipCode":0,` +
			`"LabelDate":null,"Budget":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Modification de copropriété : champ name vide"`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{Sent: []byte(`{"Copro":{"ID":0,"Reference":"CO002","Name":"Copro2",` +
			`"Address":"adresse2","ZipCode":93100,"LabelDate":"2016-04-01T12:00:00Z",` +
			`"Budget":2000000}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de copropriété, requête : Copro introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{Sent: []byte(`{"Copro":{"ID":` + strconv.Itoa(ID) + `,"Reference":"CO002",` +
			`"Name":"Copro2","Address":"adresse2","ZipCode":77001,` +
			`"LabelDate":"2016-04-01T12:00:00Z","Budget":2000000}}`),
			Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"Copro":{"ID":` + strconv.Itoa(ID) +
				`,"Reference":"CO002","Name":"Copro2","Address":"adresse2","ZipCode":77001,` +
				`"CityName":"ACHERES-LA-FORET","LabelDate":"2016-04-01T12:00:00Z","Budget":2000000`},
			StatusCode: http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/copro").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "ModifyCopro")
}

// testGetCopros check route is protected and copro are correctly sent back
func testGetCopros(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			StatusCode:   http.StatusInternalServerError}, // 0 : no token
		{Token: c.Config.Users.User.Token,
			RespContains:  []string{`"Copro"`, `"City"`, `"FcPreProg"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCopros")
}

// testGetCoproDatas check route is protected and copro datas are correctly sent back
func testGetCoproDatas(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			StatusCode:   http.StatusInternalServerError}, // 0 : no token
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Données d'une copropriété, requête copro :`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"Copro"`, `"Commitment":[],"Payment":[],` +
				`"Commission":[],"CoproForecast":[]`, `"BudgetAction"`},
			Count:         1,
			CountItemName: `"ID"`,
			ID:            ID,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro/"+strconv.Itoa(tc.ID)+"/datas").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoprosDatas")
}

// testDeleteCopro check route is protected for admin and modifications are correctly done
func testDeleteCopro(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			ID:           0,
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de copropriété, requête : Copro introuvable`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Copropriété supprimée`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/copro/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteCopro")
}

// testBatchCopros check route is limited to admin and batch import succeeds
func testBatchCopros(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Copro":[{"Reference":"","Name":"copro3","Address":"adresse3",` +
				`"ZipCode":77001,"LabelDate":null,"Budget":null},
			{"Reference":"CO004","Name":"copro4","Address":"adresse4","ZipCode":75000,` +
				`"LabelDate":42461,"Budget":3000000}]}`),
			RespContains: []string{`Batch de copropriétés, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 1 : reference empty
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Copro":[{"Reference":"","Name":"copro3","Address":"adresse3",` +
				`"ZipCode":77001,"LabelDate":null,"Budget":null},
			{"Reference":"CO004","Name":"copro4","Address":"adresse4","ZipCode":75000,` +
				`"LabelDate":42461,"Budget":3000000}]}`),
			RespContains: []string{`Batch de copropriétés, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad zip code
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Copro":[{"Reference":"CO003","Name":"copro3",` +
				`"Address":"adresse3","ZipCode":77001,"LabelDate":null,"Budget":null},
			{"Reference":"CO004","Name":"copro4","Address":"adresse4","ZipCode":75101,` +
				`"LabelDate":42461,"Budget":3000000}]}`),
			RespContains: []string{`Batch de copropriétés importé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copros").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	if chkFactory(t, tcc, f, "BatchCopro") {
		response := c.E.GET("/api/copro").
			WithHeader("Authorization", "Bearer "+c.Config.Users.Admin.Token).Expect()
		body := string(response.Content)
		for _, j := range []string{`"Reference":"CO003","Name":"copro3",` +
			`"Address":"adresse3","ZipCode":77001,"CityName":"ACHERES-LA-FORET",` +
			`"LabelDate":null,"Budget":null`,
			`"Reference":"CO004","Name":"copro4","Address":"adresse4","ZipCode":75101,` +
				`"CityName":"PARIS 1","LabelDate":"2016-04-01T00:00:00Z","Budget":3000000`} {
			if !strings.Contains(body, j) {
				t.Errorf("BatchCopro[final]\n  ->attendu %s\n  ->reçu: %s", j, body)
			}
		}
	}
}
