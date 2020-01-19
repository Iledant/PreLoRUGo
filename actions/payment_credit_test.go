package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

func testPaymentCredits(t *testing.T, c *TestContext) {
	t.Run("PaymentCredits", func(t *testing.T) {
		batchPaymentCreditsTest(t, c)
		getPaymentCreditsTest(t, c)
	})
}

// batchPaymentCreditsTest check route is admin protected and response is ok
func batchPaymentCreditsTest(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad user
		{
			Token:        c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentCredit":[`),
			RespContains: []string{"Batch d'enveloppes de crédits, décodage : "}}, // 1 : bad payload
		{
			Token:      c.Config.Users.Admin.Token,
			StatusCode: http.StatusOK,
			Sent: []byte(`{"PaymentCredit":[{"Chapter":908,"Function":811,` +
				`"Primitive":1000000,"Reported":0,"Added":500000,"Modified":300000,` +
				`"Movement":50000}]}`),
			RespContains: []string{"Enveloppes de crédits importées"}}, // 2 : ok
	}

	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/payment_credits").
			WithHeader("Authorization", "Bearer "+tc.Token).WithBytes(tc.Sent).Expect()
	}
	for _, r := range chkFactory(tcc, f, "BatchPaymentCredits") {
		t.Error(r)
	}
}

// getPaymentCreditsTest check route is protected and datas sent back are correct
func getPaymentCreditsTest(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusBadRequest,
			Params:       "a",
			RespContains: []string{`Liste des enveloppes de crédits, décodage : `}},
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusOK,
			Params:       "2019",
			RespContains: []string{`{"PaymentCredit":[`}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/payment_credits").WithHeader("Authorization", "Bearer "+tc.Token).
			WithQuery("Year", tc.Params).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetPaymentCredits") {
		t.Error(r)
	}
}
