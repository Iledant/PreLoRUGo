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

// testGetBeneficiaryDatas checks if route is user protected and datas correctly
// sent back
func testGetBeneficiaryDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=savigny`),
			RespContains: []string{`Page de données bénéficiaire, décodage Year :`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`"Datas":[],"Page":1,"ItemsCount":0`},
			Count:        0,
			ID:           0,
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{
			Token: c.Config.Users.User.Token,
			Sent:  []byte(`Page=2&Year=2010&Search=`),
			//cSpell: disable
			RespContains: []string{`"Datas":[`, `"Date":`, `"Value":`, `"Name":"`,
				`"IRISCode"`, `"Page":1`, `"ItemsCount":1`},
			//cSpell: enable
			ID:            3,
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary/"+strconv.Itoa(tc.ID)+"/datas").
			WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaryDatas") {
		t.Error(r)
	}
}

// testExportGetBeneficiaryDatas checks if route is user protected
// and datas correctly sent back
func testGetExportBeneficiaryDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Year=a&Search=savigny`),
			RespContains: []string{`Export données bénéficiaire, décodage Year :`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Year=2010&Search=savigny`),
			RespContains: []string{`"BeneficiaryData":[]`},
			Count:        0,
			ID:           0,
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{
			Token: c.Config.Users.User.Token,
			Sent:  []byte(`Year=2010&Search=`),
			//cSpell: disable
			RespContains: []string{`"BeneficiaryData":[`, `"Date"`, `"Value":`,
				`"IRISCode"`, `"Caducity":`},
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
	for _, r := range chkFactory(tcc, f, "GetExportBeneficiaryDatas") {
		t.Error(r)
	}
}
