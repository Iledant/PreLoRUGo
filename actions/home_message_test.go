package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHomeMessage is the entry point for testing all renew projet requests
func testHomeMessage(t *testing.T, c *TestContext) {
	t.Run("HomeMessage", func(t *testing.T) {
		testSetHomeMessage(t, c)
	})
}

// testSetHommeMessage checks route is admin protected and message correctly set
func testSetHomeMessage(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user not allowed
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`"Title":"`),
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Fixation du message d'accueil, décodage : `}}, // 1 : bad request
		{
			Token:      c.Config.Users.Admin.Token,
			Sent:       []byte(`{"Title":"Message du jour","Body":"Corps du message"}`),
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/home_message").
			WithHeader("Authorization", "Bearer "+tc.Token).WithBytes(tc.Sent).Expect()
	}
	chkFactory(t, tcc, f, "SetHomeMessage")
}
