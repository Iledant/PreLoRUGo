package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testReservationFee is the entry point for testing all housing comments
func testReservationFee(t *testing.T, c *TestContext) {
	t.Run("ReservationFee", func(t *testing.T) {
		ID := testCreateReservationFee(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la commentaire de logement")
			t.FailNow()
			return
		}
		testUpdateReservationFee(t, c, ID)
	})
}

// testCreateReservationFee checks if route is admin protected and created reservation
// fee is properly filled
func testCreateReservationFee(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.ReservationFeeCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Création de réservation de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent: []byte(`{"ReservationFee":{"CurrentBeneficiaryID":2,
			"PastBeneficiaryID":null,"CityCode":75101,"AddressNumber":"12",
			"AddressStreet":"rue de Vaugirard","RPLS":"RPLS123",
			"ConventionID":` + strconv.Itoa(c.HousingConventionID) + `,"Count":1,
			"TransferDate":null,"CommentID":` + strconv.Itoa(c.HousingCommentID) + `,
			"TransferID":` + strconv.Itoa(c.HousingTransferID) + `,"ConventionDate":null,
			"EliseRef":"D2020-XXXXX-00001","Area":10.23,"EndYear":2020,"Loan":350.12,
			"Charges":124.56}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"ReservationFee"`, `"ConventionDate":null`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/reservation_fee").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateReservationFee", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateReservationFee checks if route is admin protected and modified reservation
// fee is sent back
func testUpdateReservationFee(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.ReservationFeeCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de réservation de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent: []byte(`{"ReservationFee":{"ID":` + strconv.Itoa(ID) + `,"CurrentBeneficiaryID":0,
			"PastBeneficiaryID":3,"CityCode":75101,"AddressNumber":"25",
			"AddressStreet":"boulevard Pasteur","RPLS":"RPLS1234",
			"ConventionID":null,"Count":2,
			"TransferDate":"2020-01-03T00:00:00Z","CommentID":null,
			"TransferID":null,"ConventionDate":"2019-03-10T00:00:00Z",
			"EliseRef":"D2020-XXXXX-00002","Area":23.08,"EndYear":2017,"Loan":235.67,
			"Charges":99.99}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de réservation de logement, paramètre : CurrentBeneficiaryID null`},
			StatusCode:   http.StatusBadRequest}, // 2 : bad current beneficiary ID
		{
			Sent: []byte(`{"ReservationFee":{"ID":` + strconv.Itoa(ID) + `,"CurrentBeneficiaryID":2,
			"PastBeneficiaryID":3,"CityCode":75101,"AddressNumber":"25",
			"AddressStreet":"boulevard Pasteur","RPLS":"RPLS1234",
			"ConventionID":null,"Count":0,
			"TransferDate":"2020-01-03T00:00:00Z","CommentID":null,
			"TransferID":null,"ConventionDate":"2019-03-10T00:00:00Z",
			"EliseRef":"D2020-XXXXX-00002","Area":23.08,"EndYear":2017,"Loan":235.67,
			"Charges":99.99}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de réservation de logement, paramètre : Count null`},
			StatusCode:   http.StatusBadRequest}, // 3 : bad count
		{
			Sent: []byte(`{"ReservationFee":{"ID":` + strconv.Itoa(ID) + `,"CurrentBeneficiaryID":2,
			"PastBeneficiaryID":3,"CityCode":0,"AddressNumber":"25",
			"AddressStreet":"boulevard Pasteur","RPLS":"RPLS1234",
			"ConventionID":null,"Count":2,
			"TransferDate":"2020-01-03T00:00:00Z","CommentID":null,
			"TransferID":null,"ConventionDate":"2019-03-10T00:00:00Z",
			"EliseRef":"D2020-XXXXX-00002","Area":23.08,"EndYear":2017,"Loan":235.67,
			"Charges":99.99}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de réservation de logement, paramètre : CityCode null`},
			StatusCode:   http.StatusBadRequest}, // 4 : bad city code
		{
			Sent: []byte(`{"ReservationFee":{"ID":` + strconv.Itoa(ID) + `,"CurrentBeneficiaryID":2,
			"PastBeneficiaryID":123,"CityCode":75101,"AddressNumber":"25",
			"AddressStreet":"boulevard Pasteur","RPLS":"RPLS1234",
			"ConventionID":null,"Count":2,
			"TransferDate":"2020-01-03T00:00:00Z","CommentID":null,
			"TransferID":null,"ConventionDate":"2019-03-10T00:00:00Z",
			"EliseRef":"D2020-XXXXX-00002","Area":23.08,"EndYear":2017,"Loan":235.67,
			"Charges":99.99}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de réservation de logement, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 5 : past beneficiary ID doesn't
		{
			Sent: []byte(`{"ReservationFee":{"ID":0,"CurrentBeneficiaryID":2,
			"PastBeneficiaryID":3,"CityCode":75101,"AddressNumber":"25",
			"AddressStreet":"boulevard Pasteur","RPLS":"RPLS1234",
			"ConventionID":null,"Count":2,
			"TransferDate":"2020-01-03T00:00:00Z","CommentID":null,
			"TransferID":null,"ConventionDate":"2019-03-10T00:00:00Z",
			"EliseRef":"D2020-XXXXX-00002","Area":23.08,"EndYear":2017,"Loan":235.67,
			"Charges":99.99}}`),
			Token:        c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`Modification de réservation de logement, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 6 : bad ID
		{
			Sent: []byte(`{"ReservationFee":{"ID":` + strconv.Itoa(ID) + `,"CurrentBeneficiaryID":2,
			"PastBeneficiaryID":3,"CityCode":75101,"AddressNumber":"25",
			"AddressStreet":"boulevard Pasteur","RPLS":"RPLS1234",
			"ConventionID":null,"Count":2,
			"TransferDate":"2020-01-03T00:00:00Z","CommentID":null,
			"TransferID":null,"ConventionDate":"2019-03-10T00:00:00Z",
			"EliseRef":"D2020-XXXXX-00002","Area":23.08,"EndYear":2017,"Loan":235.67,
			"Charges":99.99}}`),
			Token: c.Config.Users.ReservationFeeUser.Token,
			RespContains: []string{`"ReservationFee"`, `"AddressNumber":"25"`,
				`"AddressStreet":"boulevard Pasteur"`, `"RPLS":"RPLS1234"`,
				`"ConventionID":null`, `"Count":2`,
				`"TransferDate":"2020-01-03T00:00:00Z"`, `"CommentID":null`,
				`"TransferID":null`, `"ConventionDate":"2019-03-10T00:00:00Z"`,
				`"EliseRef":"D2020-XXXXX-00002"`, `"Area":23.08`, `"EndYear":2017`,
				`"Loan":235.67`, `"Charges":99.99`},
			StatusCode: http.StatusOK}, // 7 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/reservation_fee").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateReservationFee", &ID) {
		t.Error(r)
	}
}
