package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBeneficiary is the entry point for testing all renew projet requests
func testBeneficiary(t *testing.T, c *TestContext) {
	t.Run("Beneficiary", func(t *testing.T) {
		testGetBeneficiaries(t, c)
		testGetPaginatedBeneficiaries(t, c)
	})
}

// testGetBeneficiaries checks if route is user protected and Beneficiaries correctly sent back
func testGetBeneficiaries(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"Beneficiary"`, // cSpell: disable
				`"Code":56080,"Name":"BLANGIS"`, `"Code":40505,"Name":"CLD IMMOBILIER"`,
				`"Code":7010,"Name":"IMMOBILIERE 3F"`,
				//cSpell: enable
			},
			Count:         4,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiaries").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaries") {
		t.Error(r)
	}
}

// testGetPaginatedBeneficiaries checks if route is user protected and Beneficiaries correctly sent back
func testGetPaginatedBeneficiaries(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			Sent:  []byte(`Page=2&Search=humanisme`),
			RespContains: []string{`"Beneficiary"`, `"Page"`, `"ItemsCount"`,
				// cSpell: disable
				`"Code":20186,"Name":"SCA FONCIERE HABITAT ET HUMANISME"`,
				//cSpell: enable
			},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiaries/paginated").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetPaginatedBeneficiaries") {
		t.Error(r)
	}
}
