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
	})
}

// testCreateReservationFee checks if route is admin protected and created housing
// comment is properly filled
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
