package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testSoldCommitment is the entry point for testing all renew projet requests
func testSoldCommitment(t *testing.T, c *TestContext) {
	t.Run("SoldCommitment", func(t *testing.T) {
		testGetEldestCommitments(t, c)
		testGetUnpaidCommitments(t, c)
	})
}

// testGetEldestCommitments checks if route is admin protected and SoldCommitments
// are sent back
func testGetEldestCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"SoldCommitment":[`},
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/eldest").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetEldestCommitments") {
		t.Error(r)
	}
}

// testGetUnpaidCommitments checks if route is admin protected and SoldCommitments
// are sent back
func testGetUnpaidCommitments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"SoldCommitment":[`},
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/unpaid").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetUnpaidCommitments") {
		t.Error(r)
	}
}
