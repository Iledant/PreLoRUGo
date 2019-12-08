package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCityReport is the entry point for testing all renew projet requests
func testCityReport(t *testing.T, c *TestContext) {
	t.Run("CityReport", func(t *testing.T) {
		testGetCityReport(t, c)
	})
}

// testGetCityReport checks if route is user protected and
// CityReport correctly sent back
func testGetCityReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`inseeCode=a1&firstYear=2015&lastYear=2019`),
			RespContains: []string{`Rapport par commune, décodage inseeCode :`},
			StatusCode:   http.StatusBadRequest}, // 1 : inseeCode not ok
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`inseeCode=77001&firstYear=a&lastYear=2019`),
			RespContains: []string{`Rapport par commune, décodage firstYear :`},
			StatusCode:   http.StatusBadRequest}, // 2 : firstYear not ok
		{
			Token:        c.Config.Users.User.Token,
			Sent:         []byte(`inseeCode=77001&firstYear=2015&lastYear=a`),
			RespContains: []string{`Rapport par commune, décodage lastYear :`},
			StatusCode:   http.StatusBadRequest}, // 2 : lastYear not ok
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`inseeCode=77001&firstYear=2015&lastYear=2019`),
			RespContains: []string{`"CityReport":[`,
				`{"Kind":1,"Year":2015,"Commitment":30000000,"Payment":0}`},
			StatusCode: http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/city_report").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCityReport")
}
