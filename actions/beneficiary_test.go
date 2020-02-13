package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBeneficiary is the entry point for testing all renew projet requests
func testBeneficiary(t *testing.T, c *TestContext) {
	t.Run("Beneficiary", func(t *testing.T) {
		testGetBeneficiaries(t, c)
		testGetPaginatedBeneficiaries(t, c)
		ID := testCreateBeneficiary(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le bénéficiaire")
			t.FailNow()
			return
		}
		testUpdateBeneficiary(t, c, ID)
		testDeleteBeneficiary(t, c, ID)

	})
}

// testGetBeneficiaries checks if route is user protected and Beneficiaries
// correctly sent back
func testGetBeneficiaries(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"Beneficiary"`, // cSpell: disable
				`"Code":56080,"Name":"BLANGIS"`, `"Code":40505,"Name":"CLD IMMOBILIER"`,
				`"Code":7010,"Name":"IMMOBILIERE 3F"`,
				//cSpell: enable
			},
			Count:         4,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiaries").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaries") {
		t.Error(r)
	}
}

// testGetPaginatedBeneficiaries checks if route is user protected and beneficiaries
//  correctly sent back
func testGetPaginatedBeneficiaries(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			Sent:  []byte(`Page=2&Search=humanisme`),
			RespContains: []string{`"Beneficiary"`, `"Page"`, `"ItemsCount"`,
				// cSpell: disable
				`"Code":20186,"Name":"SCA FONCIERE HABITAT ET HUMANISME"`,
				//cSpell: enable
			},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiaries/paginated").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetPaginatedBeneficiaries") {
		t.Error(r)
	}
}

// testCreateBeneficiary checks if route is admin protected and created beneficiary
// is properly filled
func testCreateBeneficiary(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de bénéficiaire, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"Beneficiary":{"Code":-1,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de bénéficiaire : name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"Beneficiary":{"Code":-1,"Name":"bénéficiaire"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"Beneficiary":`, `"Code":-1`, `"Name":"bénéficiaire"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/beneficiary").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateBeneficiary", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateBeneficiary checks if route is admin protected and updated beneficiary
// is properly filled
func testUpdateBeneficiary(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de bénéficiaire, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"Beneficiary":{"Code":-1,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de bénéficiaire : name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : code nul
		{
			Sent: []byte(`{"Beneficiary":{"ID":` + strconv.Itoa(ID) +
				`,"Code":7010,"Name":"bénéficiaire modifié"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de bénéficiaire, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 3 : duplicated beneficiary
		{
			Sent: []byte(`{"Beneficiary":{"ID":` + strconv.Itoa(ID) +
				`,"Code":-2,"Name":"bénéficiaire modifié"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Beneficiary":`, `"Code":-2`, `"Name":"bénéficiaire modifié"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/beneficiary").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateBeneficiary") {
		t.Error(r)
	}
}

// testDeleteBeneficiary checks if route is user protected and cities correctly
// sent back
func testDeleteBeneficiary(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user token
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de bénéficiaire, requête : bénéficiaire introuvable`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Bénéficiaire supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/beneficiary/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteBeneficiary") {
		t.Error(r)
	}
}
