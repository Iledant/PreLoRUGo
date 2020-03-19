package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testAvgPmtTime is the entry point for testing all average payment time requests
func testAvgPmtTime(t *testing.T, c *TestContext) {
	t.Run("Payment", func(t *testing.T) {
		testGetAvgPmtTime(t, c)
	})
}

// testGetAvgPmtTime checks if route is user protected and Payments correctly sent back
func testGetAvgPmtTime(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"AveragePaymentTime":[`},
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/avg_pmt_times").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetAvgPmtTime") {
		t.Error(r)
	}
}
