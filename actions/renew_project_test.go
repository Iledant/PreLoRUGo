package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
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
		testGetRenewProjectDatas(t, c, ID)
		testDeleteRenewProject(t, c, ID)
		testBatchRenewProject(t, c)
	})
}

// testCreateRenewProject checks if route is admin protected and created budget action
// is properly filled
func testCreateRenewProject(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"","Name":"PRU"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ reference incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"PRU001","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU","Budget":0}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : budget null
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU","Budget":250000000,"CityCode1":0}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de projet de renouvellement : Champ citycode1 incorrect`},
			StatusCode:   http.StatusBadRequest}, // 5 : CityCode1 null
		{
			Sent: []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU",` +
				`"Budget":250000000,"CityCode1":75101,"BudgetCity1":100000}}`),
			Token:  c.Config.Users.Admin.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"RenewProject":{"ID":1,"Reference":"PRU001",` +
				`"Name":"PRU","Budget":250000000,"PRIN":false,"CityCode1":75101,` +
				`"CityName1":"PARIS 1","BudgetCity1":100000,"CityCode2":null,` +
				`"CityName2":null,"BudgetCity2":null,"CityCode3":null,"CityName3":null,` +
				`"BudgetCity3":null,"Population":null,"CompositeIndex":null`},
			StatusCode: http.StatusCreated}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/renew_project").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateRenewProject", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateRenewProject checks if route is admin protected and created budget action
// is properly filled
func testUpdateRenewProject(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"","Name":"PRU"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement : Champ reference incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"PRU001","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement : Champ name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"RenewProject":{"Reference":"PRU001","Name":"PRU","Budget":0}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement : Champ budget incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : budget null
		{
			Sent:         []byte(`{"RenewProject":{"ID":0,"Reference":"PRU001","Name":"PRU","Budget":250000000,"PRIN":false,"CityCode1":75101,"CityCode2":null,"CityCode3":null,"Population":null,"CompositeIndex":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de projet de renouvellement, requête : Projet de renouvellement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{
			Sent: []byte(`{"RenewProject":{"ID":` + strconv.Itoa(ID) +
				`,"Reference":"PRU002","Name":"PRU2","Budget":150000000,"PRIN":false,` +
				`"CityCode1":77001,"CityCode2":75101,"CityCode3":78146,"Population":5400,` +
				`"CompositeIndex":1,"BudgetCity1":null,"BudgetCity2":200,"BudgetCity3":5}}`),
			Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"RenewProject":{"ID":` + strconv.Itoa(ID) +
				`,"Reference":"PRU002","Name":"PRU2","Budget":150000000,"PRIN":false,` +
				`"CityCode1":77001,"CityName1":"ACHERES-LA-FORET","BudgetCity1":null,` +
				`"CityCode2":75101,"CityName2":"PARIS 1","BudgetCity2":200,"CityCode3":78146,` +
				`"CityName3":"CHATOU","BudgetCity3":5,"Population":5400,"CompositeIndex":1}`},
			StatusCode: http.StatusCreated}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/renew_project").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateRenewProject") {
		t.Error(r)
	}
}

// testGetRenewProjects checks route is protected and all renew projects are correctly
// sent back
func testGetRenewProjects(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"RenewProject"`, `"Reference":"PRU002",` +
				`"Name":"PRU2","Budget":150000000,"PRIN":false,"CityCode1":77001,` +
				`"CityName1":"ACHERES-LA-FORET","BudgetCity1":null,"CityCode2":75101,` +
				`"CityName2":"PARIS 1","BudgetCity2":200,"CityCode3":78146,"CityName3":` +
				`"CHATOU","BudgetCity3":5,"Population":5400,"CompositeIndex":1`,
				`"City":[`, `"RPEventType":[`, `"FcPreProg":[`, `"Commission":[`,
				`"BudgetAction":[`, `"RPMultiAnnualReport":[`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : bad request
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_projects").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetRenewProjects") {
		t.Error(r)
	}
}

// testDeleteRenewProject checks that route is admin protected and delete request
// sends ok back
func testDeleteRenewProject(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.User.Token,
			ID:           0,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 1 : user unauthorized
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de projet de renouvellement, requête : Projet de renouvellement introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Projet de renouvellement supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/renew_project/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteRenewProject") {
		t.Error(r)
	}
}

// testBatchRenewProject checks that route is admin protected and batch request
// sends ok back
func testBatchRenewProject(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 1 : user unauthorized
		{
			Token: c.Config.Users.Admin.Token,
			Sent: []byte(`RenewProject":[{"Reference":"PRU002","Name":"Site RU 1","Budget":250000000},
			{"Reference":"PRU003","Name":"Site RU 2","Budget":150000000}]}`),
			RespContains: []string{`Batch de projets de renouvellement, décodage`},
			StatusCode:   http.StatusBadRequest}, // 2 : bad payload
		{
			Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"RenewProject":[{"Reference":"PRU002","Name":"Site RU 1","Budget":250000000},
			{"Reference":"PRU002","Name":"Site RU 2","Budget":150000000}]}`),
			RespContains: []string{`Batch de projets de renouvellement, requête`},
			StatusCode:   http.StatusInternalServerError}, // 3 : duplicated reference
		{
			Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"RenewProject":[{"Reference":"PRU002","Name":"Site RU 1","Budget":250000000,"PRIN":true,"CityCode1":75101,"CityCode2":null,"CityCode3":null,"Population":null,"CompositeIndex":null},
			{"Reference":"PRU003","Name":"Site RU 2","Budget":150000000,"PRIN":false,"CityCode1":77001,"CityCode2":78146,"CityCode3":null,"Population":5400,"CompositeIndex":2}]}`),
			RespContains: []string{`Batch de projets de renouvellement importé`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/renew_projects").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	resp := chkFactory(tcc, f, "BatchRenewProject")
	for _, r := range resp {
		t.Error(r)
	}
	if len(resp) > 0 {
		return
	}
	var count int64
	err := c.DB.QueryRow("SELECT count(1) FROM renew_project").Scan(&count)
	if err != nil {
		t.Errorf("Impossible de lire le nombre d'éléments insérés")
		t.FailNow()
		return
	}
	if count != 2 {
		t.Errorf("BatchRenewProject : 2 projets devaient être insérés, trouvés : %d", count)
		t.FailNow()
		return
	}
	if err = c.DB.QueryRow("SELECT id FROM renew_project WHERE reference='PRU003'").
		Scan(&c.RenewProjectID); err != nil {
		t.Errorf("Impossible de récupérer l'ID du projet de renouvellement")
		t.FailNow()
		return
	}
}

// testGetRenewProjectDatas checks that route is user protected and datas sent
// have correct fields
func testGetRenewProjectDatas(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			ID:           0,
			RespContains: []string{`Datas de projet de renouvellement, requête renewProject :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token: c.Config.Users.User.Token,
			ID:    ID,
			RespContains: []string{`"RenewProject":{"ID":` + strconv.Itoa(ID) +
				`,"Reference":"PRU002","Name":"PRU2","Budget":150000000,"PRIN":false,` +
				`"CityCode1":77001,"CityName1":"ACHERES-LA-FORET","BudgetCity1":null,` +
				`"CityCode2":75101,"CityName2":"PARIS 1","BudgetCity2":200,` +
				`"CityCode3":78146,"CityName3":"CHATOU","BudgetCity3":5,"Population":5400,` +
				`"CompositeIndex":1}`, `"Commitment"`, `"Payment"`,
				`"RenewProjectForecast"`, `"Commission"`,
				`"BudgetAction"`, `"RPEventType":[`, `"FullRPEvent":[`, `"RPCmtCityJoin":[`},
			Count:         1,
			CountItemName: "Reference",
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_project/"+strconv.Itoa(tc.ID)+"/datas").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetRenewProjectData") {
		t.Error(r)
	}
}
