package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// testUser is the entry point for testing all user related routes
func testUser(t *testing.T, c *TestContext) {
	t.Run("user", func(t *testing.T) {
		ID := testCreateUser(t, c)
		if ID == 0 {
			t.Error("Impossible de créer l'utilisateur")
			t.FailNow()
			return
		}
		testUpdateUser(t, c, ID)
		testLogout(t, c, ID)
		testChangeUserPwd(t, c)
		testGetUsers(t, c)
		testUpdateUsers(t, c, ID)
		testSetPwd(t, c, ID)
		testDeleteUser(t, c, ID)
		testSignUp(t, c)
	})
}

// testCreateUser checks route is protected and user correctly created
func testCreateUser(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Name":"essai","Email":"toto@iledefrance.fr","Password":"toto","Rights":0}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : User unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'utilisateur, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad JSON
		{Sent: []byte(`{"Name":"","Email":"toto@iledefrance.fr","Password":"toto","Rights":0}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'utilisateur : Champ name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{Sent: []byte(`{"Name":"essai","Email":"","Password":"toto","Rights":0}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'utilisateur : Champ email vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : email empty
		{Sent: []byte(`{"Name":"essai","Email":"toto@iledefrance.fr","Password":"","Rights":0}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'utilisateur : Champ password vide`},
			StatusCode:   http.StatusBadRequest}, // 4 : password empty
		{Sent: []byte(`{"Name":"Utilisateur","Email":"toto@iledefrance.fr","Password":"toto","Rights":0}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'utilisateur : Utilisateur existant`},
			StatusCode:   http.StatusBadRequest}, // 5 : name already exists
		{Sent: []byte(`{"Name":"essai","Email":"` + c.Config.Users.User.Email + `","Password":"toto","Rights":0}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'utilisateur : Utilisateur existant`},
			StatusCode:   http.StatusBadRequest}, // 6 : email already exists
		{Sent: []byte(`{"Name":"essai","Email":"toto@iledefrance.fr","Password":"toto","Rights":0}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Name":"essai","Email":"toto@iledefrance.fr","Rights":0`},
			StatusCode:   http.StatusCreated}, // 7 : correct test case
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/user").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateUser[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateUser[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"User":{"ID":%d`, &ID)
		}
	}
	return ID
}

// testLogout check if logout works
func testLogout(t *testing.T, c *TestContext, ID int) {
	response := c.E.POST("/api/user/logout").
		WithHeader("Authorization", "Bearer "+c.Config.Users.User.Token).Expect()
	body := string(response.Content)
	r := "Utilisateur déconnecté"
	if !strings.Contains(body, r) {
		t.Errorf("Logout\n  ->attendu %s\n  ->reçu: %s", r, body)
	}
	status := response.Raw().StatusCode
	if status != http.StatusOK {
		t.Errorf("Logout  ->status attendu %d  ->reçu: %d", http.StatusOK, status)
		return
	}
	response = c.E.POST("/api/user/login").WithBytes([]byte(`{"Email":"` +
		c.Config.Users.User.Email + `","Password":"` + c.Config.Users.User.Password + `"}`)).
		Expect()
	lr := struct{ Token string }{}
	if err := json.Unmarshal(response.Content, &lr); err != nil {
		t.Errorf("Logout reconnexion : " + err.Error())
		t.FailNow()
		return
	}
	c.Config.Users.User.Token = lr.Token
}

// testUpdateUser checks route is protected and user correctly modified
func testUpdateUser(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Name":"essai","Email":"toto@iledefrance.fr","Password":"toto","Rights":1}`),
			Token:        c.Config.Users.User.Token,
			ID:           ID,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`{"Name":"essai","Email":"toto@iledefrance.fr","Password":"toto","Rights":1}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Modification d'utilisateur, get`},
			StatusCode:   http.StatusInternalServerError}, // 1 : ID doesn't exist
		{Sent: []byte(`{"Name":"","Email":"toto@iledefrance.fr","Password":"toto","Rights":1}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`"Name":"essai","Email":"toto@iledefrance.fr","Rights":1`},
			StatusCode:   http.StatusOK}, // 2 : name and email unchanged
		{Sent: []byte(`{"Name":"essai2","Email":"toto2@iledefrance.fr","Rights":1}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`"Name":"essai2","Email":"toto2@iledefrance.fr","Rights":1`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/user/"+strconv.Itoa(tc.ID)).WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateUser[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateUser[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testUpdateUsers checks route is protected and user correctly modified
func testUpdateUsers(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Name":"essai","Email":"toto@iledefrance.fr","Password":"toto","Rights":1}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`{"User":[{"ID":` + strconv.Itoa(ID) + `,"Name":"","Email":"toto3@iledefrance.fr","Rights":9}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'utilisateurs, requête : Champ incorrect`},
			StatusCode:   http.StatusInternalServerError}, // 1 : name empty
		{Sent: []byte(`{"User":[{"ID":` + strconv.Itoa(ID) + `,"Name":"essai3","Email":"","Rights":9}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'utilisateurs, requête : Champ incorrect`},
			StatusCode:   http.StatusInternalServerError}, // 2 : email empty
		{Sent: []byte(`{"User":[{"ID":0,"Name":"essai3","Email":"","Rights":9}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'utilisateurs, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 3 : name and email unchanged
		{Sent: []byte(`{"User":[{"ID":` + strconv.Itoa(ID) + `,"Name":"essai3","Email":"toto3@iledefrance.fr","Rights":9}]}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`{"ID":` + strconv.Itoa(ID) + `,"Name":"essai3","Email":"toto3@iledefrance.fr","Rights":9}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/users").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateUsers[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateUsers[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testChangeUserPwd checks route is protected and user correctly modified
func testChangeUserPwd(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Sent: []byte(`{"CurrentPassword":` + c.Config.Users.User.Password + `,"NewPassword":"toto2"}`),
			Token:        "fake",
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : no token
		{Sent: []byte(`{"CurrentPassword":"","NewPassword":"toto2"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Changement de mot de passe : ancien et nouveau mots de passe requis`},
			StatusCode:   http.StatusBadRequest}, // 1 : current password empty
		{Sent: []byte(`{"CurrentPassword":"` + c.Config.Users.User.Password + `","NewPassword":""}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Changement de mot de passe : ancien et nouveau mots de passe requis`},
			StatusCode:   http.StatusBadRequest}, // 2 : new password empty
		{Sent: []byte(`{"CurrentPassword":"fake","NewPassword":"toto2"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Changement de mot de passe : erreur de mot de passe`},
			StatusCode:   http.StatusBadRequest}, // 3 : bad current password
		{Sent: []byte(`{"CurrentPassword":"` + c.Config.Users.User.Password + `","NewPassword":"toto2"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Mot de passe changé`},
			StatusCode:   http.StatusOK}, // 4 : ok
		{Sent: []byte(`{"CurrentPassword":"toto2","NewPassword":"` + c.Config.Users.User.Password + `"}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Mot de passe changé`},
			StatusCode:   http.StatusOK}, // 5 : check new password works and restore
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/user/password").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("ChangeUserPwd[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("ChangeUserPwd[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetUsers checks route is protected for admin and 3 users are sent back
func testGetUsers(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"Christophe Saintillan"`, `"essai2"`, `"Utilisateur"`},
			Count:        3,
			StatusCode:   http.StatusOK}, // 0 : user unauthorized
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/users").WithHeader("Authorization", "Bearer "+tc.Token).
			Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetUsers[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetUsers[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.Count != 0 {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetUsers[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testSetPwd checks route is protected and user correctly modified
func testSetPwd(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Password":"mdp"}`),
			Token:        c.Config.Users.User.Token,
			ID:           ID,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`{"Password":"mdp"}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Modification de mot de passe, get`},
			StatusCode:   http.StatusInternalServerError}, // 1 : ID doesn't exist
		{Sent: []byte(`{"Password":""}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Modification de mot de passe, mot de passe vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : password empty
		{Sent: []byte(`{"Password":"mdp"}`),
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`"User":{"ID":3,"Name":"essai3","Email":"toto3@iledefrance.fr","Rights":9}`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/users/pwd/"+strconv.Itoa(tc.ID)).WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("SetPwd[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("SetPwd[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testDeleteUser checks route is protected and user correctly modified
func testDeleteUser(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			ID:           ID,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression d'utilisateur, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Utilisateur supprimé`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.DELETE("/api/user/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("DeleteUser[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("DeleteUser[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testSignUp checks a the user is created and inactive
func testSignUp(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Sent: []byte(`fake}`),
			RespContains: []string{`Inscription d'utilisateur, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 0 : user unauthorized
		{Sent: []byte(`{"Name":""}`),
			RespContains: []string{`Inscription d'utilisateur : Champ manquant ou incorrect`},
			StatusCode:   http.StatusBadRequest}, // 1 : name empty
		{Sent: []byte(`{"Name":"Utilisateur","Email":""}`),
			RespContains: []string{`Inscription d'utilisateur : Champ manquant ou incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : email empty
		{Sent: []byte(`{"Name":"Utilisateur","Email":"user@iledefrance.fr","Password":""}`),
			RespContains: []string{`Inscription d'utilisateur : Champ manquant ou incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : password empty
		{Sent: []byte(`{"Name":"Utilisateur","Email":"user@iledefrance.fr","Password":"tutu"}`),
			RespContains: []string{`Inscription d'utilisateur, exists :`},
			StatusCode:   http.StatusBadRequest}, // 4 : users exists
		{Sent: []byte(`{"Name":"Utilisateur2","Email":"user2@iledefrance.fr","Password":"tutu"}`),
			RespContains: []string{`Utilisateur créé, en attente d'activation`},
			StatusCode:   http.StatusCreated}, // 5 : created
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/user/sign_up").WithBytes(tc.Sent).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("SignUp[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("SignUp[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}
