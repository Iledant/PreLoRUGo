package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

func testPaymentCreditJournals(t *testing.T, c *TestContext) {
	t.Run("PaymentCreditJournals", func(t *testing.T) {
		batchPaymentCreditJournalsTest(t, c)
		getPaymentCreditJournalsTest(t, c)
	})
}

// batchPaymentCreditJournalsTest check route is admin protected and response is ok
func batchPaymentCreditJournalsTest(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase,
		{
			Token:        c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentCredit":[`),
			RespContains: []string{"Batch mouvements de crédits, décodage : "}},
		{
			Token:      c.Config.Users.Admin.Token,
			StatusCode: http.StatusOK,
			Sent: []byte(`{"PaymentCreditJournal":[{"Chapter":908,"Function":811,` +
				`"CreationDate":20190310,"ModificationDate":20190315,"Name":"Mouvement","Value":100000}]}`),
			RespContains: []string{"Mouvements de crédits importés"}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/payment_credit_journal").
			WithHeader("Authorization", "Bearer "+tc.Token).WithBytes(tc.Sent).Expect()
	}
	chkFactory(t, tcc, f, "BatchPaymentCreditJournals")
}

// getPaymentCreditJournalsTest check route is protected and datas sent back are correct
func getPaymentCreditJournalsTest(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusBadRequest,
			Params:       "a",
			RespContains: []string{`Mouvements de crédits, décodage : `}},
		{
			Token:      c.Config.Users.User.Token,
			StatusCode: http.StatusOK,
			Params:     "2019",
			RespContains: []string{`{"PaymentCreditJournal":[{"Chapter":908,"ID":1,` +
				`"Function":811,"CreationDate":"2019-03-10T00:00:00Z","ModificationDate"` +
				`:"2019-03-15T00:00:00Z","Name":"Mouvement","Value":100000}]}`}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/payment_credit_journal").
			WithHeader("Authorization", "Bearer "+tc.Token).WithQuery("Year", tc.Params).Expect()
	}
	chkFactory(t, tcc, f, "GetPaymentCreditJournals")
}
