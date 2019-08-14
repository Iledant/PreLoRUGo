package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetBeneficiaryPayments handle the get request to fetch the payment per month
// of a beneficiary
func GetBeneficiaryPayments(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Paiement d'un bénéficiaire, erreur ID : " + err.Error()})
		return
	}
	var resp models.BeneficiaryPayments
	db := ctx.Values().Get("db").(*sql.DB)
	if err = resp.GetAll(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Paiement d'un bénéficiaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
