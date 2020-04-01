package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

func testPaymentDemands(t *testing.T, c *TestContext) {
	t.Run("PaymentDemands", func(t *testing.T) {
		batchPaymentDemandsTest(t, c)
		updatePaymentDemandsTest(t, c)
		getAllPaymentDemandsTest(t, c)
		getPaymentDemandCountsTest(t, c)
	})
}

// batchPaymentDemandsTest check route is protected and a small batch doesn't raise error
func batchPaymentDemandsTest(t *testing.T, c *TestContext) {
	TestCases := []TestCase{
		*c.AdminCheckTestCase, // 0 unauthorized
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}]}`),
			RespContains: []string{"Batch de demandes de paiement, décodage"}}, // 1 bad json
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 IrisCode vide"}}, // 2 IrisCode empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 IrisName vide"}}, //3 IrisName empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 CommitmentDate vide"}}, // 4 commiment_date empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 BeneficiaryCode vide"}}, // 5 BeneficiaryCode empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 DemandNumber vide"}}, // 6 demande_number empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 DemandDate vide"}}, // 7 DemandDate empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demandes de paiement : ligne 1 ReceiptDate vide"}}, // 8 ReceiptDate empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusBadRequest,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}]}`),
			RespContains: []string{"Batch de demandes de paiement : date d'import non définie"}}, // 9 import_date empty
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusOK,
			Sent:         []byte(`{"PaymentDemand":[{"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":43168,"BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":43268,"ReceiptDate":43278,"DemandValue":10000000,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null}],"ImportDate":"2018-06-28T01:00:00Z"}`),
			RespContains: []string{"Batch de demande de paiement importé"}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/payment_demands").
			WithHeader("Authorization", "Bearer "+tc.Token).WithBytes(tc.Sent).Expect()
	}
	for _, r := range chkFactory(TestCases, f, "BatchPaymentDemands") {
		t.Error(r)
	}
}

// updatePaymentDemandsTest check route is protected and payment demands correctly
// sent back
func updatePaymentDemandsTest(t *testing.T, c *TestContext) {
	TestCases := []TestCase{
		*c.AdminCheckTestCase,
		{Token: c.Config.Users.Admin.Token,
			StatusCode:   http.StatusOK,
			Sent:         []byte(`{"PaymentDemand":{"id":1,"IrisCode":"12000139","IrisName":"Construction logements PLAI","CommitmentDate":"2018-03-09T01:00:00Z","BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":"2018-06-17T01:00:00Z","ReceiptDate":"2018-06-27T01:00:00Z","CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null,"IrisCode":"12000139","IrisName":"Construction logements PLAI","BeneficiaryCode":7010,"DemandNumber":1,"DemandDate":"2018-06-17T01:00:00Z","ReceiptDate":"2018-06-27T01:00:00Z","DemandValue":545600,"CsfDate":null,"CsfComment":null,"DemandStatus":null,"StatusComment":null,"Excluded":true,"ExcludedComment":"commentaire"}}`),
			RespContains: []string{`"ID":1`, `"commentaire"`, `"Excluded":true`, ``}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/payment_demands").
			WithHeader("Authorization", "Bearer "+tc.Token).WithBytes(tc.Sent).Expect()
	}
	for _, r := range chkFactory(TestCases, f, "UpdatePaymentDemands") {
		t.Error(r)
	}
}

// getAllPaymentDemandsTest check route is protected and payment demands correctly sent.
func getAllPaymentDemandsTest(t *testing.T, c *TestContext) {
	TestCases := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:         c.Config.Users.User.Token,
			StatusCode:    http.StatusOK,
			RespContains:  []string{`"PaymentDemand"`, `"IrisName":"Construction logements PLAI"`},
			CountItemName: `"ID"`,
			Count:         1},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/payment_demands").WithQueryString("Param").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(TestCases, f, "GetAllPaymentDemands") {
		t.Error(r)
	}
}

// getPaymentDemandCountsTest check route is protected and payment demands correctly sent.
func getPaymentDemandCountsTest(t *testing.T, c *TestContext) {
	TestCases := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:         c.Config.Users.User.Token,
			StatusCode:    http.StatusOK,
			RespContains:  []string{`"PaymentDemandCount"`, `"Unprocessed"`, `"Uncontrolled"`},
			CountItemName: `"Date"`,
			Count:         31},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/payment_demand_counts").WithQueryString("Param").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(TestCases, f, "GetPaymentDemandCounts") {
		t.Error(r)
	}
}
