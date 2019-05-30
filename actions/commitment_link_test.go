package actions

import (
	"net/http"
	"testing"
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
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{`),
			RespContains: []string{"Liens d'engagements, décodage : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad payload
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":0,"IDs":[1],"Type":"Copro"}`),
			RespContains: []string{"Liens d'engagements, format : ID d'engagement incorrect"},
			StatusCode:   http.StatusBadRequest}, // 2 : DestID nul
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":1,"IDs":[1],"Type":""}`),
			RespContains: []string{"Liens d'engagements, format : Type incorrect"},
			StatusCode:   http.StatusBadRequest}, // 3 : bad type
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":1,"IDs":[1],"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad copro ID
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":2,"IDs":[2,3,5],"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad commitment ID
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":3,"IDs":[2,3],"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements mis à jour`},
			StatusCode:   http.StatusOK}, // 6 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/commitments/link").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		chkBodyStatusAndCount(t, tc, i, response, "SetCommitmentLink")
	}
}

// testUnlinkCommitment check route is limited to admin and batch import succeeds
func testUnlinkCommitment(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{`),
			RespContains: []string{"Suppression de liens d'engagements, décodage : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad payload
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IDs":[2,3,5]}`),
			RespContains: []string{`Suppression de liens d'engagements, requête : Impossible de supprimer tous les liens`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad commitment ID
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IDs":[2,3]}`),
			RespContains: []string{`Suppression de liens d'engagements mis à jour`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/commitments/unlink").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		chkBodyStatusAndCount(t, tc, i, response, "UnlinkCommitmentLink")
	}
}
