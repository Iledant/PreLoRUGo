package actions

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// testPayment is the entry point for testing all renew projet requests
func testPayment(t *testing.T, c *TestContext) {
	t.Run("Payment", func(t *testing.T) {
		testBatchPayments(t, c)
		testGetPayments(t, c)
	})
}

// testBatchPayments check route is limited to admin and batch import succeeds
func testBatchPayments(t *testing.T, c *TestContext) {
	batchContent, err := ioutil.ReadFile("../assets/payment_batch.json")
	if err != nil {
		t.Errorf("Impossible de lire le ficher de batch")
		return
	}
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"Payment":[{"CommitmentYear":2012,"CommitmentCode":"IRIS "}]}`),
			RespContains: []string{"Batch de Paiements, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : validation error
		{Token: c.Config.Users.Admin.Token,
			Sent:         batchContent,
			RespContains: []string{"Batch de Paiements importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/payments").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchPayment[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchPayment[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		// GelAllTest checks if data are correctly analyzed
	}
}

// testGetPayments checks if route is user protected and Payments correctly sent back
func testGetPayments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"CommitmentYear":2012,"CommitmentCode":"IRIS ","CommitmentNumber":392543,"CommitmentLine":1,"Year":2014,"CreationDate":"2014-01-29T00:00:00Z","ModificationDate":"2014-02-10T00:00:00Z","Value":12648324`},
			Count:        3,
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/payments").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetPayments[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetPayments[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetPayments[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
