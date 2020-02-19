package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetDifActionPaymentPrevisions handle the get request to calculate the payment
// previsions per action using the past commitments, the programmation of the
// actual year and the housing, copro and renew project forecast for the coming years
func GetDifActionPaymentPrevisions(ctx iris.Context) {
	var resp models.DifActionPmtPrevisions
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Prévisions de paiement par action, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
