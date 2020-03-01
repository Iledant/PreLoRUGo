package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetEldestCommitments handle the get request to fetch the eldest commitments
// in order to checks which ones should be sold out
func GetEldestCommitments(ctx iris.Context) {
	var resp models.SoldCommitments
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetOld(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des engagements anciens, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetUnpaidCommitments handle the get request to fetch the eldest commitments
// in order to checks which ones should be sold out
func GetUnpaidCommitments(ctx iris.Context) {
	var resp models.SoldCommitments
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetUnpaid(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des engagements non payés, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
