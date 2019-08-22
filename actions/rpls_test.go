package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPLS is the entry point for testing all renew projet requests
func testRPLS(t *testing.T, c *TestContext) {
	t.Run("RPLS", func(t *testing.T) {
		ID := testCreateRPLS(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le RPLS")
			t.FailNow()
			return
		}
		testUpdateRPLS(t, c, ID)
		testDeleteRPLS(t, c, ID)
		testBatchRPLS(t, c)
		testGetAllRPLS(t, c)
		testGetRPLSDatas(t, c)
		testRPLSReport(t, c)
		testRPLSDetailedReport(t, c)
	})
}

// testCreateRPLS checks if route is admin protected and created RPLS
// is properly filled
func testCreateRPLS(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"RPLS":{"InseeCode":75101,"Year":2016,"Ratio":0.167}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de RPLS, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"Year":2016,"Ratio":0.167}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Création de RPLS, requête : InseeCode nul"`},
			StatusCode:   http.StatusInternalServerError}, // 2 : InseeCode nul
		{
			Sent:         []byte(`{"InseeCode":75101,"Ratio":0.167}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Création de RPLS, requête : Year nul"`},
			StatusCode:   http.StatusInternalServerError}, // 3 : Year nul
		{
			Sent:         []byte(`{"InseeCode":75000,"Year":2016,"Ratio":0.167}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de RPLS, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad InseeCode
		{
			Sent:   []byte(`{"InseeCode":75101,"Year":2016,"Ratio":0.167}`),
			Token:  c.Config.Users.Admin.Token,
			IDName: `"ID"`,
			RespContains: []string{`{"RPLS":{"ID":2,"InseeCode":75101,"CityName":"PARIS 1",` +
				`"Year":2016,"Ratio":0.167}}`},
			StatusCode: http.StatusCreated}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/rpls").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateRPLS", &ID)
	return ID
}

// testUpdateRPLS checks if route is admin protected and RPLS
// is properly modified
func testUpdateRPLS(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"RPLS":{"InseeCode":75101,"Year":2016,"Ratio":0.167}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de RPLS, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"InseeCode":77001,"Year":2017,"Ratio":0.3}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de RPLS, requête : ID nul`},
			StatusCode:   http.StatusInternalServerError}, // 2 : ID nul
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"Year":2017,"Ratio":0.3}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de RPLS, requête : InseeCode nul`},
			StatusCode:   http.StatusInternalServerError}, // 3 : InseeCode nul
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"InseeCode":77001,"Ratio":0.3}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de RPLS, requête : Year nul`},
			StatusCode:   http.StatusInternalServerError}, // 4 : Year nul
		{
			Sent: []byte(`{"ID":` + strconv.Itoa(ID) + `,"InseeCode":77000,` +
				`"Year":2017,"Ratio":0.3}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de RPLS, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad InseeCode
		{
			Sent: []byte(`{"ID":` + strconv.Itoa(ID) +
				`,"InseeCode":77001,"Year":2017,"Ratio":0.3}`),
			Token: c.Config.Users.Admin.Token,
			RespContains: []string{`{"RPLS":{"ID":` + strconv.Itoa(ID) +
				`,"InseeCode":77001,"CityName":"ACHERES-LA-FORET","Year":2017,"Ratio":0.3}`},
			StatusCode: http.StatusOK}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/rpls").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateRPLS")
}

// testDeleteRPLS checks if route is admin protected and RPLS
// is properly modified
func testDeleteRPLS(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			ID:           0,
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			ID:           0,
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de RPLS, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			ID:           ID,
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`RPLS supprimé`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/rpls/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteRPLS")
}

// testBatchRPLS checks if route is admin protected and created RPLS
// is properly filled
func testBatchRPLS(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"RPLS":[{"InseeCode":75101,"Year":2016,"Ratio":0.167}]`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Batch RPLS, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent: []byte(`{"RPLS":[{"Year":2016,"Ratio":0.167},` +
				`{"InseeCode":77101,"Year":2016,"Ratio":0.3},` +
				`{"InseeCode":78146,"Year":2016,"Ratio":0.1955}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Batch RPLS, requête : ligne 1 InseeCode null"`},
			StatusCode:   http.StatusInternalServerError}, // 2 : InseeCode nul
		{
			Sent: []byte(`{"RPLS":[{"InseeCode":75101,"Year":2016,"Ratio":0.167},` +
				`{"InseeCode":77101,"Ratio":0.3},` +
				`{"InseeCode":78146,"Year":2016,"Ratio":0.1955}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Batch RPLS, requête : ligne 2 Year null"`},
			StatusCode:   http.StatusInternalServerError}, // 3 : Year nul
		{
			Sent: []byte(`{"RPLS":[{"InseeCode":75101,"Year":2016,"Ratio":0.167},` +
				`{"InseeCode":77001,"Year":2016,"Ratio":0.3},` +
				`{"InseeCode":78146,"Year":2016,"Ratio":0.1955}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Batch RPLS importé"`},
			StatusCode:   http.StatusOK}, // 4 : ok
		{
			Sent:         []byte(`{"RPLS":[{"InseeCode":78146,"Year":2016,"Ratio":0.4}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Batch RPLS importé"`},
			StatusCode:   http.StatusOK}, // 5 : ok, test updating
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/rpls/batch").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchRPLS")
}

// testGetAllRPLS check if route is user protected and batch RPLS have been
// correctly inserted and updated
func testGetAllRPLS(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: `fake`,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Token invalid`}},
		{Token: c.Config.Users.User.Token,
			StatusCode: http.StatusOK,
			RespContains: []string{`{"RPLS":[{"ID":3,"InseeCode":75101,"CityName":` +
				`"PARIS 1","Year":2016,"Ratio":0.167},{"ID":4,"InseeCode":77001,` +
				`"CityName":"ACHERES-LA-FORET","Year":2016,"Ratio":0.3},{"ID":5,` +
				`"InseeCode":78146,"CityName":"CHATOU","Year":2016,"Ratio":0.4}]}`}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rpls").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetAllRPLS")

}

// testGetRPLSDatas check if route is admin protected and batch datas are
// correctly sent back
func testGetRPLSDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: `fake`,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Token invalid`}},
		{Token: c.Config.Users.Admin.Token,
			StatusCode: http.StatusOK,
			RespContains: []string{`"RPLS":[{"ID":3,"InseeCode":75101,"CityName":` +
				`"PARIS 1","Year":2016,"Ratio":0.167},{"ID":4,"InseeCode":77001,` +
				`"CityName":"ACHERES-LA-FORET","Year":2016,"Ratio":0.3},{"ID":5,` +
				`"InseeCode":78146,"CityName":"CHATOU","Year":2016,"Ratio":0.4}]`,
				`"City":[`}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rpls/datas").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPLSDatas")

}

// testRPLSReport check if route is user protected and batch RPLS have been
// correctly inserted and updated
func testRPLSReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{
			Token:        `fake`,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Token invalid`}}, // 0 bad token
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=a&FirstYear=2010&LastYear=2019&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport RPLS, décodage RPLSYear :`}}, // 1 bad RPLSYear
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=a&LastYear=2019&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport RPLS, décodage FirstYear :`}}, // 2 bad FirstYear
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=a&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport RPLS, décodage LastYear :`}}, // 3 bad LastYear
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=2019&RPLSMin=a&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport RPLS, décodage RPLSMin :`}}, // 4 bad RPLSMin
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=2019&RPLSMin=0&RPLSMax=a",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport RPLS, décodage RPLSMax :`}}, // 5 bad RPLSMax
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=2019&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusOK,
			RespContains: []string{`{"RPLSReport":[`}}, // 6 ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rpls/report").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "RPLSReport")

}

// testRPLSDetailedReport check if route is user protected and batch RPLS have been
// correctly inserted and updated
func testRPLSDetailedReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{
			Token:        `fake`,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Token invalid`}}, // 0 bad token
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=a&FirstYear=2010&LastYear=2019&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport détaillé RPLS, décodage RPLSYear :`}}, // 1 bad RPLSYear
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=a&LastYear=2019&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport détaillé RPLS, décodage FirstYear :`}}, // 2 bad FirstYear
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=a&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport détaillé RPLS, décodage LastYear :`}}, // 3 bad LastYear
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=2019&RPLSMin=a&RPLSMax=0.3",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport détaillé RPLS, décodage RPLSMin :`}}, // 4 bad RPLSMin
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=2019&RPLSMin=0&RPLSMax=a",
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Rapport détaillé RPLS, décodage RPLSMax :`}}, // 5 bad RPLSMax
		{
			Token:        c.Config.Users.User.Token,
			Params:       "RPLSYear=2016&FirstYear=2010&LastYear=2019&RPLSMin=0&RPLSMax=0.3",
			StatusCode:   http.StatusOK,
			RespContains: []string{`{"RPLSDetailedReport":[`}}, // 6 ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rpls/detailed_report").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "RPLSDetailedReport")

}
