package actions

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// testCommitment is the entry point for testing all renew projet requests
func testCommitment(t *testing.T, c *TestContext) {
	t.Run("Commitment", func(t *testing.T) {
		testBatchCommitments(t, c)
		testGetCommitments(t, c)
	})
}

// testBatchCommitments check route is limited to admin and batch import succeeds
func testBatchCommitments(t *testing.T, c *TestContext) {
	correctBatch, err := ioutil.ReadFile("../assets/commitment_batch.json")
	if err != nil {
		t.Errorf("Impossible de lire le fichier commitment_batch.json")
	}
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"Commitment":[{"Year":2010}]}`),
			RespContains: []string{"Batch de Engagements, requête : Champs incorrects"},
			StatusCode:   http.StatusInternalServerError}, // 1 : validation error
		{Token: c.Config.Users.Admin.Token,
			Sent:         correctBatch,
			RespContains: []string{"Batch de Engagements importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/commitments").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchCommitment[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchCommitment[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		// the testGetCommitments is used to check datas have benne correctly imported
	}
}

// testGetCommitments checks if route is user protected and Commitments correctly sent back
func testGetCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			// cSpell: disable
			RespContains: []string{`"Commitment"`, `"Year":2009,"Code":"AE   ","Number":244923,"Line":1,"CreationDate":"2012-01-26T00:00:00Z","ModificationDate":"2012-01-26T00:00:00Z","Name":"TRAITEMENT DE CADUCITE 2011","Value":-15371500,"BeneficiaryID":3,"IrisCode":null`},
			// cSpell: enable
			Count:      4,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/commitments").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetCommitments[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetCommitments[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetCommitments[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
