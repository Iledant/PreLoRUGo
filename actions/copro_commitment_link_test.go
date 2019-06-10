package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCoproCommitmentLink is the entry point for testing settings requests
func testCoproCommitmentLink(t *testing.T, c *TestContext) {
	t.Run("CoproCommitmentLink", func(t *testing.T) {
		testLinkCommitmentsCopros(t, c)
	})
}

// testLinkCommitmentsCopros checks route is protected and all settings are
// correctly sent back
func testLinkCommitmentsCopros(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les copropriétés requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBach":[{"Reference":"","IRISCode":"13021233"}]}`),
			RespContains: []string{`Ligne 1, Reference vide`},
			StatusCode:   http.StatusInternalServerError}, // 1 : reference null
		{Token: c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBach":[{"Reference":"Essai3","IRISCode":""}]}`),
			RespContains: []string{`Ligne 1, IRISCode vide`},
			StatusCode:   http.StatusInternalServerError}, // 2 : IRIS code empty
		{Token: c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBach":[{"Reference":"Essai3","IRISCode":"1302123"}]}`),
			RespContains: []string{`Liens engagements copros, requête : ligne 0 Reference ou code IRIS introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad reference
		{Token: c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBach":[{"Reference":"CO004","IRISCode":"13021233"}]}`),
			RespContains: []string{`Liens engagements copros importés`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro/commitments").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "LinkCommitmentsCopros")
}
