package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBeneficiaryPayments is the entry point for testing all renew projet requests
func testBeneficiaryPayments(t *testing.T, c *TestContext) {
	t.Run("BeneficiaryPayments", func(t *testing.T) {
		testGetBeneficiaryPayments(t, c)
	})
}

// testGetBeneficiaryPayments checks if route is user protected and datas correctly sent back
func testGetBeneficiaryPayments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"BeneficiaryPayment":[]`},
			Count:        0,
			ID:           0,
			StatusCode:   http.StatusOK}, // 1 : bad ID
		{
			Token: c.Config.Users.User.Token,
			//cSpell: disable
			RespContains: []string{`BeneficiaryPayment":[{"Year":2010,"Month":2,` +
				`"Value":18968.8},{"Year":2011,"Month":6,"Value":4742.2}]`},
			//cSpell: enable
			ID:            4,
			Count:         2,
			CountItemName: `"Year"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary/"+strconv.Itoa(tc.ID)+"/payments").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaryPayments") {
		t.Error(r)
	}
}
