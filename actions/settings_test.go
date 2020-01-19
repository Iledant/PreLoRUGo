package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
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
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetSector":[`, `"BudgetAction":[`, `"Commission":[`,
				`"PaginatedCity":{`, `"Community":[`, `"PaginatedPayment":{`, `"PaginatedCommitment":{`},
			Count:         15,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : bad request
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/settings").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetSettings") {
		t.Error(r)
	}
}
