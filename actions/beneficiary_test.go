package actions

import (
	"net/http"
	"strings"
	"testing"
)

// testBeneficiary is the entry point for testing all renew projet requests
func testBeneficiary(t *testing.T, c *TestContext) {
	t.Run("Beneficiary", func(t *testing.T) {
		testGetBeneficiaries(t, c)
	})
}

// testGetBeneficiaries checks if route is user protected and Beneficiaries correctly sent back
func testGetBeneficiaries(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"Beneficiary"`, `"Code":5019,"Name":"ESSONNE HABITAT                       "`,
				// cSpell: disable
				`"Code":30953,"Name":"FONCIA MARCEAU                        "`,
				`"Code":6850,"Name":"SA D HLM LOGIREP                      "`,
				//cSpell: enable
				`"Code":29364,"Name":"OPH MANTES YVELINES HABITAT           "`},
			Count:      4,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/beneficiaries").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetBeneficiaries[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetBeneficiaries[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetBeneficiaries[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
