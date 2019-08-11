package actions

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCommitment is the entry point for testing all renew projet requests
func testCommitment(t *testing.T, c *TestContext) {
	t.Run("Commitment", func(t *testing.T) {
		testBatchCommitments(t, c)
		testGetCommitments(t, c)
		testGetPaginatedCommitments(t, c)
		testGetUnlinkedCommitments(t, c)
		testExportedCommitments(t, c)
	})
}

// testBatchCommitments check route is limited to admin and batch import succeeds
func testBatchCommitments(t *testing.T, c *TestContext) {
	correctBatch, err := ioutil.ReadFile("../assets/commitment_batch.json")
	if err != nil {
		t.Errorf("Impossible de lire le fichier commitment_batch.json")
	}
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"Commitment":[{"Year":2010}]}`),
			RespContains: []string{"Batch de Engagements, requête : Ligne 1 : champs incorrects"},
			StatusCode:   http.StatusInternalServerError}, // 1 : validation error
		{Token: c.Config.Users.Admin.Token,
			Sent:         correctBatch,
			RespContains: []string{"Batch de Engagements importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/commitments").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchCommitment")
	// the testGetCommitments is used to check datas have been correctly imported
}

// testGetCommitments checks if route is user protected and Commitments correctly sent back
func testGetCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			// cSpell: disable
			RespContains: []string{`"Commitment"`, `"Year":2012,"Code":"IRIS","Number"` +
				`:362012,"Line":1,"CreationDate":"2012-02-02T00:00:00Z","ModificationDate":` +
				`"2012-02-02T00:00:00Z","CaducityDate":"2015-05-02T00:00:00Z",` +
				`"Name":"12000139 - 1","Value":232828,"SoldOut":true,"BeneficiaryID":1,` +
				`"ActionID":2,"IrisCode":"12000139","HousingID":null,"CoproID":null,` +
				`"RenewProjectID":null`},
			// cSpell: enable
			Count:         4,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCommitments")
}

// testGetPaginatedCommitments checks if route is user protected and paginated
// commitments correctly sent back
func testGetPaginatedCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=savigny`),
			RespContains: []string{`Page d'engagements, décodage Year :`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad params query
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Page=2&Year=2010&Search=savigny`),
			// cSpell: disable
			RespContains: []string{`"Commitment"`, `"Year":2015,"Code":"IRIS ",` +
				`"Number":469347,"Line":1,"CreationDate":"2015-04-13T00:00:00Z",` +
				`"ModificationDate":"2015-04-13T00:00:00Z","CaducityDate":` +
				`"2018-06-13T00:00:00Z","Name":"91 - SAVIGNY SUR ORGE - AV DE LONGJUMEAU` +
				` - 65 PLUS/PLAI","Value":30000000,"SoldOut":false,"BeneficiaryID":3,` +
				`"BeneficiaryName":"IMMOBILIERE 3F","ActionName":"Aide à la création de ` +
				`logements locatifs très sociaux","Sector":"LO","IrisCode":"14004240",` +
				`"HousingID":null,"CoproID":null,"RenewProjectID":null`,
				`"Page":1`, `"ItemsCount":1`},
			// cSpell: enable
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/paginated").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetPaginatedCommitments")
}

// testGetUnlinkedCommitments checks if route is user protected and paginated
// commitments correctly sent back
func testGetUnlinkedCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=savigny`),
			RespContains: []string{`Page d'engagements non liés, décodage Year :`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad params query
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Page=2&Year=2010&Search=savigny`),
			// cSpell: disable
			RespContains: []string{`"Commitment"`, `"Year":2015,"Code":"IRIS ",` +
				`"Number":469347,"Line":1,"CreationDate":"2015-04-13T00:00:00Z",` +
				`"ModificationDate":"2015-04-13T00:00:00Z","CaducityDate":` +
				`"2018-06-13T00:00:00Z","Name":"91 - SAVIGNY SUR ORGE - AV DE LONGJUMEAU` +
				` - 65 PLUS/PLAI","Value":30000000,"SoldOut":false,"BeneficiaryID":3,` +
				`"BeneficiaryName":"IMMOBILIERE 3F","ActionName":` +
				`"Aide à la création de logements locatifs très sociaux","Sector":"LO",` +
				`"IrisCode":"14004240","HousingID":null,"CoproID":null,"RenewProjectID"` +
				`:null`, `"Page":1`, `"ItemsCount":1`},
			// cSpell: enable
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/unlinked").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetUnlinkedCommitments")
}

// testExportedCommitments checks if route is user protected and exported
// commitments correctly sent back
func testExportedCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Year=2010&Search=savigny`),
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Year=a&Search=savigny`),
			RespContains: []string{`Export d'engagements, décodage Year :`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad params query
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Year=2010&Search=savigny`),
			// cSpell: disable
			RespContains: []string{`"ExportedCommitment":[`, `"ID":3,"Year":2015,` +
				`"Code":"IRIS ","Number":469347,"Line":1,"CreationDate":` +
				`"2015-04-13T00:00:00Z","ModificationDate":"2015-04-13T00:00:00Z",` +
				`"CaducityDate":"2018-06-13T00:00:00Z","Name":"91 - SAVIGNY SUR ORGE - ` +
				`AV DE LONGJUMEAU - 65 PLUS/PLAI","Value":300000,"SoldOut":false,` +
				`"BeneficiaryName":"IMMOBILIERE 3F","Sector":"LO","ActionName":` +
				`"Aide à la création de logements locatifs très sociaux","IrisCode":` +
				`"14004240","HousingName":null,"CoproName":null,"RenewProjectName":null`},
			// cSpell: enable
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/export").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetExportedCommitments")
}
