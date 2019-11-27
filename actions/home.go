package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// homeResp embeddes the data sent back to the home get request
type homeResp struct {
	Commitment models.TwoYearsCommitments `json:"Commitment"`
	Payment    models.TwoYearsPayments    `json:"Payment"`
	models.CumulatedProgrammation
	models.ImportLogs
	models.PaymentCreditSum
}

// GetHome handle the get request for the home page
func GetHome(ctx iris.Context) {
	var resp homeResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Commitment.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Home requête engagement : " + err.Error()})
		return
	}
	if err := resp.Payment.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Home requête paiement : " + err.Error()})
		return
	}
	if err := resp.CumulatedProgrammation.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Home requête programmation : " + err.Error()})
		return
	}
	if err := resp.ImportLogs.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Home requête logs : " + err.Error()})
		return
	}
	if err := resp.PaymentCreditSum.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Home requête payment credit sum : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
