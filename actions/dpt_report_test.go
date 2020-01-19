package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testDepartmentReport is the entry point for testing all renew projet requests
func testDepartmentReport(t *testing.T, c *TestContext) {
	t.Run("DepartmentReport", func(t *testing.T) {
		testGetDepartmentReport(t, c)
	})
}

// testGetDepartmentReport checks if route is user protected and
// DepartmentReport correctly sent back
func testGetDepartmentReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`firstYear=a&lastYear=2019`),
			RespContains: []string{`Rapport par département, décodage firstYear :`},
			StatusCode:   http.StatusBadRequest}, // 1 : firstYear not ok
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`firstYear=2016&lastYear=a`),
			RespContains: []string{`Rapport par département, décodage lastYear :`},
			StatusCode:   http.StatusBadRequest}, // 2 : lastYear not ok
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`firstYear=2016&lastYear=2019`),
			RespContains: []string{`"DptReport":[]`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/department_report").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetDepartmentReport") {
		t.Error(r)
	}
}
