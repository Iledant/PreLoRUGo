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
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBatch":[{"Reference":"","IRISCode":"13021233"}]}`),
			RespContains: []string{`Ligne 1, Reference vide`},
			StatusCode:   http.StatusInternalServerError}, // 1 : reference null
		{
			Token:        c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBatch":[{"Reference":"Essai3","IRISCode":""}]}`),
			RespContains: []string{`Ligne 1, IRISCode vide`},
			StatusCode:   http.StatusInternalServerError}, // 2 : IRIS code empty
		{
			Token:        c.Config.Users.CoproUser.Token,
			Sent:         []byte(`{"CoproCommitmentBatch":[{"Reference":"CO004","IRISCode":"13021233"}]}`),
			RespContains: []string{`Liens engagements copros importés`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro/commitments").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	resp := chkFactory(tcc, f, "LinkCommitmentsCopros")
	for _, r := range resp {
		t.Error(r)
	}
	if len(resp) > 0 {
		return
	}
	var count int64
	if err := c.DB.QueryRow(`SELECT count(1) FROM commitment 
		WHERE iris_code='13021233' AND copro_id IS NOT NULL`).Scan(&count); err != nil {
		t.Errorf("LinkCommitmentsCopros, erreur sur la requête de vérification %v", err)
	}
	if count != 1 {
		t.Errorf("LinkCommitmentsCopros: échec de la vérification -> attendu 1  -> reçu %d\n", count)
	}
}
