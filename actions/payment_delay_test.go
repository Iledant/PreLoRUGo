package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

func testPaymentDelays(t *testing.T, c *TestContext) {
	t.Run("PaymentDelays", func(t *testing.T) {
		getPaymentDelaysTest(t, c)
	})
}

// getPaymentDelaysTest check route is protected and datas sent back are correct
func getPaymentDelaysTest(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusBadRequest,
			Params:       "a",
			RespContains: []string{`Délais de paiement, décodage : `}},
		{
			Token:         c.Config.Users.User.Token,
			StatusCode:    http.StatusOK,
			Params:        "1472688000000",
			RespContains:  []string{`{"PaymentDelay":[`, `"Number":1`},
			Count:         12,
			CountItemName: `"Delay"`},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/payment_delays").WithHeader("Authorization", "Bearer "+tc.Token).
			WithQuery("after", tc.Params).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetPaymentDelays") {
		t.Error(r)
	}
}
