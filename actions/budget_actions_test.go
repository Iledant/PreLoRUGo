package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// testBudgetAction is the entry point for testing all budget action related routes
func testBudgetAction(t *testing.T, c *TestContext) {
	t.Run("BudgetAction", func(t *testing.T) {
		ID := testCreateBudgetAction(t, c)
		if ID == 0 {
			t.Error("Impossible de créer l'action budgétaire")
			t.FailNow()
			return
		}
		testUpdateBudgetAction(t, c, ID)
		testGetBudgetActions(t, c)
		testDeleteBudgetAction(t, c, ID)
	})
}

// testCreateBudgetAction checks if route is admin protected and created budget action
// is properly filled
func testCreateBudgetAction(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"BudgetAction":{"Code":"1234567890","Name":"Action"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'action budgétaire, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : code empty
		{Sent: []byte(`{"BudgetAction":{"Code":"","Name":"Action"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code empty
		{Sent: []byte(`{"BudgetAction":{"Code":"1234567890","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"BudgetAction":{"Code":"1234567890","Name":"Action"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetAction"`, `"Code":"1234567890","Name":"Action"`},
			StatusCode:   http.StatusCreated}, // 4 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/budget_action").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateBudgetAction[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateBudgetAction[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"BudgetAction":{"ID":%d`, &ID)
		}
	}
	return ID
}

// testGetBudgetActions checks route is protected and datas sent back are well formed
func testGetBudgetActions(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			StatusCode:   http.StatusInternalServerError}, // 0 : token null
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`BudgetAction`, `Name`, `Code`},
			Count:        1,
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/budget_actions").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetBudgetActions[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetBudgetActions[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetBudgetActions[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testUpdateBudgetAction checks if route is admin protected and budget action
// is properly modified
func testUpdateBudgetAction(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"BudgetAction":{"Code":"1234567890","Name":"Action"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : code empty
		{Sent: []byte(`{"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":"","Name":"Action"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code empty
		{Sent: []byte(`{"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":"1234567890","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"BudgetAction":{"ID":0,"Code":"1234567890","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'action budgétaire : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : name empty
		{Sent: []byte(`{"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":"0123456789","Name":"Action modifiée"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"BudgetAction":{"ID":` + strconv.Itoa(ID) + `,"Code":"0123456789","Name":"Action modifiée","SectorID":null}}`},
			StatusCode:   http.StatusOK}, // 5 : name empty
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/budget_action").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateBudgetAction[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateBudgetAction[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testDeleteBudgetAction check route is admin protected and delete requests returns ok
func testDeleteBudgetAction(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression d'action budgétaire, requête : Action budgétaire introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Action budgétaire supprimée`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.DELETE("/api/budget_action/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("DeleteBudgetAction[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("DeleteBudgetAction[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}
