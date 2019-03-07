package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type beneficiaryReq struct {
	Beneficiary models.Beneficiary `json:"Beneficiary"`
}

// GetBeneficiaries handles the get request to fetch all beneficiaries
func GetBeneficiaries(ctx iris.Context) {
	var resp models.Beneficiaries
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
