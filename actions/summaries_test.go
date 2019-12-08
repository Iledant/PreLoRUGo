package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testSummaries is the entry point for testing all renew projet requests
func testSummaries(t *testing.T, c *TestContext) {
	t.Run("Summaries", func(t *testing.T) {
		testGetSummariesDatas(t, c)
	})
}

// testGetSummariesDatas check if route is user protected and summaries datas are
// correctly sent back
func testGetSummariesDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase,
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusOK,
			RespContains: []string{`"City":[`, `"RPLSYear":[2016]`}},
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/summaries/datas").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetSummariesDatas")

}
