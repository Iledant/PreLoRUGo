package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testReservationReport is the entry point for testing all reservation reports
func testReservationReport(t *testing.T, c *TestContext) {
	t.Run("ReservationReport", func(t *testing.T) {
		ID := testCreateReservationReport(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le report de réservation")
			t.FailNow()
			return
		}
		testUpdateReservationReport(t, c, ID)
		testGetReservationReports(t, c)
		testDeleteReservationReport(t, c, ID)
	})
}

// testCreateReservationReport checks if route is reservation protected and
// created reservation report is properly filled
func testCreateReservationReport(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.ReservationFeeCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Création de report de réservation, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"ReservationReport":{"BeneficiaryID":0,"Area":10.5,"SourceIRISCode":"EX0001","DestIRISCode":null,"DestDate":null}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Création de report de réservation, paramètre : BeneficiaryID vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : BeneficiaryID null
		{
			Sent:         []byte(`{"ReservationReport":{"BeneficiaryID":0,"Area":10.5,"DestIRISCode":null,"DestDate":null}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Création de report de réservation, paramètre : BeneficiaryID vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : SourceIRISCode vide
		{
			Sent:         []byte(`{"ReservationReport":{"BeneficiaryID":2,"Area":10.5,"SourceIRISCode":"EX0001","DestIRISCode":null,"DestDate":null}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"ReservationReport"`, `"BeneficiaryID":2`, `"Area":10.5`, `"SourceIRISCode":"EX0001"`, `"DestIRISCode":null`, `"DestDate":null`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/reservation_report").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateReservationReport", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateReservationReport checks if route is reservation protected and
// reservation report is properly filled
func testUpdateReservationReport(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.ReservationFeeCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de report de réservation, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"ReservationReport":{"ID":0,"BeneficiaryID":3,"Area":15.5,"SourceIRISCode":"EX0002","DestIRISCode":"EX111","DestDate":"2019-03-01T00:00:00Z"}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de report de réservation, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"ReservationReport":{"ID":` + strconv.Itoa(ID) + `,"BeneficiaryID":0,"Area":15.5,"SourceIRISCode":"EX0002","DestIRISCode":"EX111","DestDate":"2019-03-01T00:00:00Z"}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de report de réservation, paramètre : BeneficiaryID vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : BeneficiaryID empty
		{
			Sent:         []byte(`{"ReservationReport":{"ID":` + strconv.Itoa(ID) + `,"BeneficiaryID":3,"Area":15.5,"DestIRISCode":"EX111","DestDate":"2019-03-01T00:00:00Z"}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de report de réservation, paramètre : SourceIRISCode vide`},
			StatusCode:   http.StatusBadRequest}, // 4 : SourceIRISCode empty
		{
			Sent:         []byte(`{"ReservationReport":{"ID":` + strconv.Itoa(ID) + `,"BeneficiaryID":3,"Area":15.5,"SourceIRISCode":"EX0002","DestIRISCode":"EX111","DestDate":"2019-03-01T00:00:00Z"}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`"ReservationReport"`, `"Beneficiary":"CLD IMMOBILIER"`, `"BeneficiaryID":3`, `"Area":15.5`, `"SourceIRISCode":"EX0002"`, `"DestIRISCode":"EX111"`, `"DestDate":"2019-03-01T00:00:00Z"`},
			StatusCode:   http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/reservation_report").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateReservationReport") {
		t.Error(r)
	}
}

// testGetReservationReports checks route is protected and all housing conventions
// are correctly sent back
func testGetReservationReports(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.ReservationFeeCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.ReservationFeeUser.Token,
			RespContains:  []string{`"ReservationReport"`, `"Beneficiary":"CLD IMMOBILIER"`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/reservation_reports").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetReservationReports") {
		t.Error(r)
	}
}

// testDeleteReservationReport checks that route is reservation protected and
// delete request sends ok back
func testDeleteReservationReport(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.ReservationFeeCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.ReservationFeeUser.Token,
			ID:           0,
			RespContains: []string{`Suppression de report de réservation, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.ReservationFeeUser.Token,
			ID:           ID,
			RespContains: []string{`Report de réservation supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/reservation_report/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteReservationReport") {
		t.Error(r)
	}
}
