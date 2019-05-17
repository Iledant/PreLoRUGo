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
		testGetPaginatedPayments(t, c)
		testExportedPayments(t, c)
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
			RespContains: []string{`"ID":1,"CommitmentID":4,"CommitmentYear":2010,"CommitmentCode":"IRIS ","CommitmentNumber":277678,"CommitmentLine":1,"Year":2010,"CreationDate":"2010-02-02T00:00:00Z","ModificationDate":"2010-04-16T00:00:00Z","Number":102717,"Value":1896880`,
				`"ID":4,"CommitmentID":2,"CommitmentYear":2014,"CommitmentCode":"IRIS ","CommitmentNumber":431370,"CommitmentLine":1,"Year":2016,"CreationDate":"2016-09-12T00:00:00Z","ModificationDate":"2016-09-19T00:00:00Z","Number":141103,"Value":239200`},
			Count:      4,
			StatusCode: http.StatusOK}, // 1 : ok
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

// testGetPaginatedPayments checks if route is user protected and Payments correctly sent back
func testGetPaginatedPayments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Page=2&Year=2010&Search=cld`),
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=cld`),
			RespContains: []string{`Page de paiements, décodage Year :`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Page=2&Year=2010&Search=cld`),
			//cSpell: disable
			RespContains: []string{`"Payments":[`, `"Year":2016,"CreationDate":"2016-09-12T00:00:00Z","Value":239200,"Number":141103,"CommitmentDate":"2014-02-05T00:00:00Z","CommitmentName":"13021233 - 1","CommitmentValue":239200,"Beneficiary":"CLD IMMOBILIER","Sector":"LO","ActionName":"Aide aux copropriétés en difficulté"`},
			//cSpell: enable
			Count:      1,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/payments/paginated").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetPaginatedPayments[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetPaginatedPayments[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetPaginatedPayments[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testExportedPayments checks if route is user protected and Payments correctly sent back
func testExportedPayments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Year=2010&Search=cld`),
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Year=a&Search=cld`),
			RespContains: []string{`Export de paiements, décodage Year :`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Year=2010&Search=cld`),
			//cSpell: disable
			RespContains: []string{`"ExportedPayment":[`, `"Year":2016,"CreationDate":"2016-09-12T00:00:00Z","ModificationDate":"2016-09-19T00:00:00Z","Number":141103,"Value":2392,"CommitmentYear":2014,"CommitmentCode":"IRIS ","CommitmentNumber":431370,"CommitmentCreationDate":"2014-02-05T00:00:00Z","CommitmentValue":2392,"CommitmentName":"13021233 - 1","BeneficiaryName":"CLD IMMOBILIER","Sector":"LO","ActionName":"Aide aux copropriétés en difficulté"`},
			//cSpell: enable
			Count:      1,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/payments/export").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("ExportedPayments[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("ExportedPayments[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("ExportedPayments[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
