package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

func testCommitmentLink(t *testing.T, c *TestContext) {
	t.Run("CommitmentLink", func(t *testing.T) {
		testLinkCommitment(t, c)
		testUnlinkCommitment(t, c)
	})
}

// testLinkCommitment check route is limited to admin and batch import succeeds
func testLinkCommitment(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{`),
			RespContains: []string{"Liens d'engagements, décodage : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad payload
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":0,"IDs":[1],"Type":"Copro"}`),
			RespContains: []string{"Liens d'engagements, format : ID d'engagement incorrect"},
			StatusCode:   http.StatusBadRequest}, // 2 : DestID nul
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":1,"IDs":[1],"Type":""}`),
			RespContains: []string{"Liens d'engagements, format : Type incorrect"},
			StatusCode:   http.StatusBadRequest}, // 3 : bad type
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":1,"IDs":[1],"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad copro ID
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":2,"IDs":[2,3,5],"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad commitment ID
		{
			Token:         c.Config.Users.Admin.Token,
			Sent:          []byte(`{"DestID":` + strconv.FormatInt(c.CoproID, 10) + `,"IDs":[2,3],"Type":"Copro"}`),
			RespContains:  []string{`"Commitment":[{`, `"Payment":[{`},
			Count:         3,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 6 : copro, ok
		{
			Token:         c.Config.Users.Admin.Token,
			Sent:          []byte(`{"DestID":2,"IDs":[1],"Type":"RenewProject"}`),
			RespContains:  []string{`"Commitment":[{`, `"Payment":[`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 7 : renew project, ok
		{
			Token:         c.Config.Users.Admin.Token,
			Sent:          []byte(`{"DestID":3,"IDs":[4],"Type":"Housing"}`),
			RespContains:  []string{`"Commitment":[{`, `"Payment":[`},
			Count:         3,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 7 : housing, ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/commitments/link").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "LinkCommitment") {
		t.Error(r)
	}
}

// testUnlinkCommitment check route is limited to admin and batch import succeeds
func testUnlinkCommitment(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{`),
			RespContains: []string{"Suppression de liens d'engagements, décodage : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad payload
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IDs":[2,3,5]}`),
			RespContains: []string{`Suppression de liens d'engagements, requête : Impossible de supprimer tous les liens`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad commitment ID
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IDs":[2,3]}`),
			RespContains: []string{`Suppression de liens d'engagements mis à jour`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/commitments/unlink").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UnlinkCommitmentLink") {
		t.Error(r)
	}
}
