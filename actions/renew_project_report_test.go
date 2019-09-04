package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRenewProjectReport is the entry point for testing all renew projet requests
func testRenewProjectReport(t *testing.T, c *TestContext) {
	t.Run("RenewProjectReport", func(t *testing.T) {
		testGetRenewProjectReport(t, c)
	})
}

// testGetRenewProjectReport checks if route is user protected and
// RenewProjectReport correctly sent back
func testGetRenewProjectReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"RenewProjectReport":[{"ID":2,"Reference":"PRU002",` +
				`"Name":"Site RU 1","Budget":250000000,"Commitment":null,"Payment":null,` +
				`"LastEventName":null,"LastEventDate":null,"City1Name":"PARIS 1","City1Cmt"` +
				`:null,"City1Pmt":null,"City2Name":null,"City2Cmt":null,"City2Pmt":null,` +
				`"City3Name":null,"City3Cmt":null,"City3Pmt":null},{"ID":3,"Reference":` +
				`"PRU003","Name":"Site RU 2","Budget":150000000,"Commitment":null,"Payment":` +
				`null,"LastEventName":null,"LastEventDate":null,"City1Name":"ACHERES-LA-FORET",` +
				`"City1Cmt":null,"City1Pmt":null,"City2Name":"CHATOU","City2Cmt":null,` +
				`"City2Pmt":null,"City3Name":null,"City3Cmt":null,"City3Pmt":null}]`},
			StatusCode: http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_project/report").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRenewProjectReport")
}
