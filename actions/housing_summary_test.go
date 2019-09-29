package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingSummary is the entry point for testing housing summary requests
func testHousingSummary(t *testing.T, c *TestContext) {
	t.Run("HousingSummary", func(t *testing.T) {
		testBatchHousingSummary(t, c)
	})
}

var (
	tcSentBH1 = []byte(`{`)
	tcSentBH2 = []byte(`{"HousingSummary":[{"InseeCode":,"Address":,"PLAI":,"PLS":,"PLUS":,"IRISCode":,"ReferenceCode":,"ANRU":}]}`)
)

// testBatchHousingSummary check route is limited to housing user and batch
// import succeeds
func testBatchHousingSummary(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits logement requis`},
			StatusCode:   http.StatusUnauthorized,
		}, // 0 : token empty
		{
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Batch de bilan logements, décodage :`},
			Sent:         tcSentBH1,
			StatusCode:   http.StatusBadRequest,
		}, // 1 : token empty
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_summary").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchHousingSummary")
}
