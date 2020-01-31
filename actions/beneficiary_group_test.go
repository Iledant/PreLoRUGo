package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testBeneficiaryGroup is the entry point for testing all beneficiary groups
func testBeneficiaryGroup(t *testing.T, c *TestContext) {
	t.Run("BeneficiaryGroup", func(t *testing.T) {
		ID := testCreateBeneficiaryGroup(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le groupe de bénéficiaires")
			t.FailNow()
			return
		}
		testUpdateBeneficiaryGroup(t, c, ID)
		testGetBeneficiaryGroups(t, c)
		testSetBeneficiaryGroup(t, c, ID)
		testGetBeneficiaryGroupItems(t, c, ID)
		testDeleteBeneficiaryGroup(t, c, ID)
	})
}

// testCreateBeneficiaryGroup checks if route is admin protected and created
// beneficiary group is properly filled
func testCreateBeneficiaryGroup(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de groupe de bénéficiaires, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"BeneficiaryGroup":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de groupe de bénéficiaires, paramètre : name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"BeneficiaryGroup":{"Name":"Groupe"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"BeneficiaryGroup"`, `"Name":"Groupe"`},
			StatusCode:   http.StatusCreated}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/beneficiary_group").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateBeneficiaryGroup", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{
		{
			Sent:       []byte(`{"BeneficiaryGroup":{"Name":"Groupe de test"}}`),
			Token:      c.Config.Users.Admin.Token,
			IDName:     `"ID"`,
			StatusCode: http.StatusCreated}, // 4 : ok
	}
	for _, r := range chkFactory(tcc, f, "CreateBeneficiaryGroupID", &c.BeneficiaryGroupID) {
		t.Error(r)
	}
	return ID
}

// testUpdateBeneficiaryGroup checks if route is admin protected and modified
//  beneficiary group is properly sent back
func testUpdateBeneficiaryGroup(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de groupe de bénéficiaires, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"BeneficiaryGroup":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de groupe de bénéficiaires, paramètre : name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"BeneficiaryGroup":{"ID":` + strconv.Itoa(ID) + `,"Name":"Groupe modifié"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"BeneficiaryGroup"`, `"Name":"Groupe modifié"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/beneficiary_group").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateBeneficiaryGroup") {
		t.Error(r)
	}
}

// testGetBeneficiaryGroups checks if route is user protected and modified
// group if properly sent back
func testGetBeneficiaryGroups(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"BeneficiaryGroup":[`, `"Name":"Groupe modifié"`},
			CountItemName: `"ID"`,
			Count:         2,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary_groups").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaryGroups") {
		t.Error(r)
	}
}

// testSetBeneficiaryGroup checks if route is admin protected and created
// beneficiary group is properly filled
func testSetBeneficiaryGroup(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Fixation de groupe de bénéficiaires, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"BeneficiaryIDs":[1,2,4]}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Fixation de groupe de bénéficiaires, requête`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad group ID
		{
			Sent:         []byte(`{"BeneficiaryIDs":[1,2,1000]}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Fixation de groupe de bénéficiaires, requête`},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad beneficiary ID
		{
			Sent:          []byte(`{"BeneficiaryIDs":[1,2,4]}`),
			Token:         c.Config.Users.Admin.Token,
			ID:            ID,
			CountItemName: `"ID"`,
			Count:         3,
			RespContains:  []string{`"Beneficiary":[`, `"Name":"SCA FONCIERE HABITAT ET HUMANISME"`},
			StatusCode:    http.StatusOK}, // 4 : ok
		{
			Sent:       []byte(`{"BeneficiaryIDs":[1,2,4]}`),
			Token:      c.Config.Users.Admin.Token,
			ID:         c.BeneficiaryGroupID,
			StatusCode: http.StatusOK}, // 5 : fill beneficiary test group
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/beneficiary_group/"+strconv.Itoa(tc.ID)).WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "SetBeneficiaryGroup", &ID) {
		t.Error(r)
	}
}

// testGetBeneficiaryGroupItems checks if route is user protected and modified
// group if properly sent back
func testGetBeneficiaryGroupItems(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"Beneficiary":[]`},
			ID:           0,
			StatusCode:   http.StatusOK}, // 1 : bad ID
		{
			Token: c.Config.Users.User.Token,
			// cSpell:disable
			RespContains: []string{`"Beneficiary":[`, `"Name":"BLANGIS"`, `"Code":7010`,
				`"Name":"SCA FONCIERE HABITAT ET HUMANISME"`},
			//cSpell:enable
			ID:            ID,
			CountItemName: `"ID"`,
			Count:         3,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary_group/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaryGroupItems") {
		t.Error(r)
	}
}

// testDeleteBeneficiaryGroup checks if route is admin protected and delete
// sends ok back
func testDeleteBeneficiaryGroup(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de groupe de bénéficiaires, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Groupe de bénéficiaires supprimé`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/beneficiary_group/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteBeneficiaryGroup") {
		t.Error(r)
	}
}
