package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingCommitmentLink is the entry point for testing settings requests
func testHousingCommitmentLink(t *testing.T, c *TestContext) {
	t.Run("HousingCommitmentLink", func(t *testing.T) {
		testLinkCommitmentsHousings(t, c)
	})
}

// testLinkCommitmentsHousings checks route is protected and all settings are
// correctly sent back
func testLinkCommitmentsHousings(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.HousingCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         []byte(`{"HousingCommitmentBach":[{"Reference":"","IRISCode":"14004240"}]}`),
			RespContains: []string{`Ligne 1, Reference vide`},
			StatusCode:   http.StatusInternalServerError}, // 1 : reference null
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         []byte(`{"HousingCommitmentBach":[{"Reference":"Essai3","IRISCode":""}]}`),
			RespContains: []string{`Ligne 1, IRISCode vide`},
			StatusCode:   http.StatusInternalServerError}, // 2 : IRISCode null
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         []byte(`{"HousingCommitmentBach":[{"Reference":"Essai3","IRISCode":"14004240"}]}`),
			RespContains: []string{`Liens engagements logements importés`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing/commitments").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "LinkCommitmentsHousings") {
		t.Error(r)
	}
}
