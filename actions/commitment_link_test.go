package actions

import (
	"net/http"
	"strings"
	"testing"
)

func testCommitmentLink(t *testing.T, c *TestContext) {
	t.Run("CommitmentLink", func(t *testing.T) {
		testSetCommitmentLink(t, c)
	})
}

// testSetCommitmentLink check route is limited to admin and batch import succeeds
func testSetCommitmentLink(t *testing.T, c *TestContext) {
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
			Sent:         []byte(`{"DestID":0,"IDs":[1],"Link":true,"Type":"Copro"}`),
			RespContains: []string{"Liens d'engagements, format : ID d'engagement incorrect"},
			StatusCode:   http.StatusBadRequest}, // 2 : DestID nul
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":1,"IDs":[1],"Link":true,"Type":""}`),
			RespContains: []string{"Liens d'engagements, format : Type incorrect"},
			StatusCode:   http.StatusBadRequest}, // 3 : bad type
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":1,"IDs":[1],"Link":true,"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad copro ID
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":2,"IDs":[2,3,5],"Link":true,"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements, requête : Impossible de modifier tous les ID d'engagement`},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad commitment ID
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"DestID":2,"IDs":[2,3],"Link":true,"Type":"Copro"}`),
			RespContains: []string{`Liens d'engagements mis à jour`},
			StatusCode:   http.StatusOK}, // 6 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/commitments/link").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("SetCommitmentLink[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("SetCommitmentLink[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}
