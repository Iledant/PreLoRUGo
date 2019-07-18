package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBeneficiaryDatas is the entry point for testing all renew projet requests
func testBeneficiaryDatas(t *testing.T, c *TestContext) {
	t.Run("BeneficiaryDatas", func(t *testing.T) {
		testGetBeneficiaryDatas(t, c)
		testGetExportBeneficiaryDatas(t, c)
	})
}

// testGetBeneficiaryDatas checks if route is user protected and datas correctly sent back
func testGetBeneficiaryDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`Token absent`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=savigny`),
			RespContains: []string{`Page de données bénéficiaire, décodage Year :`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`"Datas":[],"Page":1,"ItemsCount":0`},
			Count:        0,
			ID:           0,
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Page=2&Year=2010&Search=savigny`),
			//cSpell: disable
			RespContains: []string{`"Datas":[{"ID":3,"Date":"2015-04-13T00:00:00Z",` +
				`"Value":30000000,"Name":"91 - SAVIGNY SUR ORGE - AV DE LONGJUMEAU - 65 PLUS/PLAI",` +
				`"IRISCode":"14004240","Available":30000000}],"Page":1,"ItemsCount":1`},
			//cSpell: enable
			ID:            3,
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary/"+strconv.Itoa(tc.ID)+"/datas").
			WithQueryString(string(tc.Sent)).WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetBeneficiaryDatas")
}

// testExportGetBeneficiaryDatas checks if route is user protected
// and datas correctly sent back
func testGetExportBeneficiaryDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Year=2010&Search=savigny`),
			RespContains: []string{`Token absent`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Year=a&Search=savigny`),
			RespContains: []string{`Export données bénéficiaire, décodage Year :`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Year=2010&Search=savigny`),
			RespContains: []string{`"BeneficiaryData":[]`},
			Count:        0,
			ID:           0,
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Year=2010&Search=savigny`),
			//cSpell: disable
			RespContains: []string{`"BeneficiaryData":[{"ID":3,"Date":"2015-04-13T00:00:00Z",` +
				`"Value":30000000,"Name":"91 - SAVIGNY SUR ORGE - AV DE LONGJUMEAU - 65 PLUS/PLAI",` +
				`"IRISCode":"14004240","Available":30000000}]`},
			//cSpell: enable
			ID:            3,
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary/"+strconv.Itoa(tc.ID)+"/export").
			WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetExportBeneficiaryDatas")
}
