package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHome is the entry point for testing all home requests
func testHome(t *testing.T, c *TestContext) {
	t.Run("Home", func(t *testing.T) {
		testGetHome(t, c)
	})
}

// testGetHome checks if route is user protected and datas correctly sent back
func testGetHome(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains:  []string{`"Commitment":`, `"Payment":`},
			Count:         0,
			CountItemName: `"Month"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/home").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetHomes")
}
