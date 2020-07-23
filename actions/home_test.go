package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHome is the entry point for testing all home requests
func testHome(t *testing.T, c *TestContext) {
	t.Run("Home", func(t *testing.T) {
		testGetHome(t, c)
	})
}

// testGetHome checks if route is user protected and datas correctly sent back
func testGetHome(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"Commitment":`, `"Payment":`, `"ImportLog":[`,
				`"Programmation":[`, `"PaymentCreditSum":`,
				`"HomeMessage":{"Title":"Message du jour","Body":"Corps du message"}`,
				`"PaymentDemandsStock"`, `"AveragePayment":[`, `"CsfWeekTrend":`,
				`"FlowStockDelays":`, `"PaymentRate":`},
			Count:         4,
			CountItemName: `"Month"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/home").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetHomes") {
		t.Error(r)
	}
}
