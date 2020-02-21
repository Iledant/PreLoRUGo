package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testDifActionPaymentPrevisions is the entry point for testing all renew projet requests
func testDifActionPaymentPrevisions(t *testing.T, c *TestContext) {
	t.Run("DifActionPaymentPrevisions", func(t *testing.T) {
		testGetDifActionPaymentPrevisions(t, c)
	})
}

// testGetDifActionPaymentPrevisions check if route is user protected and summaries datas are
// correctly sent back
func testGetDifActionPaymentPrevisions(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusOK,
			RespContains: []string{`"DifActionPmtPrevision":[`}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/dif_action_pmt_prev").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetDifActionPaymentPrevisions") {
		t.Error(r)
	}

}
