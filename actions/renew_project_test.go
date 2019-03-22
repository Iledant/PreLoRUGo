package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Iledant/PreLoRUGo/models"
)

// testRenewProject is the entry point for testing all renew projet requests
func testRenewProject(t *testing.T, c *TestContext) {
	t.Run("RenewProject", func(t *testing.T) {
		ID := testCreateRenewProject(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le projet de renouvellement")
			t.FailNow()
			return
		}
		testUpdateRenewProject(t, c, ID)
		testGetRenewProjects(t, c)
		testDeleteRenewProject(t, c, ID)
		rp := models.RenewProject{Reference: "RP_TEST",
			Name:           "Projet RU test",
			Budget:         250000000,
			Population:     models.NullInt64{Valid: false},
			CompositeIndex: models.NullInt64{Valid: false}}
		if err := rp.Create(c.DB); err != nil {
			t.Error("Impossible de créer le projet de renouvellement de test")
			t.FailNow()
			return
		}
		c.RenewProjectID = rp.ID
	})
}

// testCreateRenewProject checks if route is admin protected and created budget action
// is properly filled
func testCreateRenewProject(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"RenewProject":{"Code":"PRU001","Name":"PRU"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"RenewProject":{"Reference":"","Name":"PRU"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ reference, name ou budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{Sent: []byte(`{"RenewProject":{"Reference":"PRU001","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ reference, name ou budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU","Budget":0}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ reference, name ou budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : budget null
		{Sent: []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU","Budget":250000000}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"RenewProject":{"ID":1,"Reference":"PRU001","Name":"PRU","Budget":250000000,"Population":null,"CompositeIndex":null}`},
			StatusCode:   http.StatusCreated}, // 5 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/renew_project").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateRenewProject[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateRenewProject[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"RenewProject":{"ID":%d`, &ID)
		}
	}
	return ID
}

// testUpdateRenewProject checks if route is admin protected and created budget action
// is properly filled
func testUpdateRenewProject(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"RenewProject":{"Code":"PRU001","Name":"PRU"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"RenewProject":{"Reference":"","Name":"PRU"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement : Champ reference, name ou budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{Sent: []byte(`{"RenewProject":{"Reference":"PRU001","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement : Champ reference, name ou budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU","Budget":0}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement : Champ reference, name ou budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : budget null
		{Sent: []byte(`{"RenewProject":{"ID":0,"Reference":"PRU002","Name":"PRU2","Budget":150000000,"Population":5400,"CompositeIndex":1}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement, requête : Projet de renouvellement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{Sent: []byte(`{"RenewProject":{"ID":` + strconv.Itoa(ID) + `,"Reference":"PRU002","Name":"PRU2","Budget":150000000,"Population":5400,"CompositeIndex":1}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"RenewProject":{"ID":1,"Reference":"PRU002","Name":"PRU2","Budget":150000000,"Population":5400,"CompositeIndex":1}`},
			StatusCode:   http.StatusCreated}, // 6 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/renew_project").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateRenewProject[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateRenewProject[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetRenewProjects checks route is protected and all renew projects are correctly
// sent back
func testGetRenewProjects(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "fake",
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"RenewProject"`, `"Name":"PRU2"`},
			Count:        1,
			StatusCode:   http.StatusOK}, // 1 : bad request
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/renew_projects").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetRenewProjects[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetRenewProjects[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetRenewProjects[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testDeleteRenewProject checks that route is admin protected and delete request
// sends ok back
func testDeleteRenewProject(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "fake",
			ID:           0,
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : bad token
		{Token: c.Config.Users.User.Token,
			ID:           0,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 1 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de projet de renouvellement, requête : Projet de renouvellement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{Token: c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Projet de renouvellement supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.DELETE("/api/renew_project/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("DeleteRenewProject[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("DeleteRenewProject[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}
