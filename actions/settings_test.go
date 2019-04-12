package actions

import (
	"net/http"
	"strings"
	"testing"
)

// testSettings is the entry point for testing settings requests
func testSettings(t *testing.T, c *TestContext) {
	t.Run("Settings", func(t *testing.T) {
		testGetSettings(t, c)
	})
}

// testGetSettings checks route is protected and all settings are correctly
// sent back
func testGetSettings(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetSector":[`, `"BudgetAction":[`, `"Commission":[`, `"City":[`, `"Community":[`},
			Count:        9,
			StatusCode:   http.StatusOK}, // 1 : bad request
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/settings").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetSettings[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetSettings[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetSettings[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
