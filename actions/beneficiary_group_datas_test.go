package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBeneficiaryGroupDatas is the entry point for testing all renew projet requests
func testBeneficiaryGroupDatas(t *testing.T, c *TestContext) {
	t.Run("BeneficiaryGroupDatas", func(t *testing.T) {
		testGetBeneficiaryGroupDatas(t, c)
		testGetExportBeneficiaryGroupDatas(t, c)
	})
}

// testGetBeneficiaryGroupDatas checks if route is user protected and datas correctly
// sent back
func testGetBeneficiaryGroupDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=savigny`),
			RespContains: []string{`Page de données groupe de bénéficiaires, décodage Year :`},
			ID:           c.BeneficiaryGroupID,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`"Datas":[],"Page":1,"ItemsCount":0`},
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{
			Token: c.Config.Users.User.Token,
			Sent:  []byte(`Page=2&Year=2010&Search=`),
			RespContains: []string{`"Datas":[`, `"Date":`, `"Value":`, `"Name":"`,
				`"IRISCode"`, `"Page":1`, `"BeneficiaryName":`, `"ItemsCount":3`},
			ID:            c.BeneficiaryGroupID,
			Count:         3,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary_group/"+strconv.Itoa(tc.ID)+"/datas").
			WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaryGroupDatas") {
		t.Error(r)
	}
}

// testExportGetBeneficiaryGroupDatas checks if route is user protected
// and datas correctly sent back
func testGetExportBeneficiaryGroupDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Year=a&Search=savigny`),
			RespContains: []string{`Export données groupe de bénéficiaires, décodage Year :`},
			ID:           c.BeneficiaryGroupID,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`Year=2010&Search=savigny`),
			RespContains: []string{`"BeneficiaryGroupData":[]`},
			ID:           0,
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{
			Token: c.Config.Users.User.Token,
			Sent:  []byte(`Year=2010&Search=`),
			//cSpell: disable
			RespContains: []string{`"BeneficiaryGroupData":[`, `"Date"`, `"Value":`,
				`"BeneficiaryName":`, `"IRISCode"`, `"Caducity":`},
			//cSpell: enable
			ID:            c.BeneficiaryGroupID,
			Count:         3,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary_group/"+strconv.Itoa(tc.ID)+"/export").
			WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetExportBeneficiaryGroupDatas") {
		t.Error(r)
	}
}
