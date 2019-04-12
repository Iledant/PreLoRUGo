package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type settingsResp struct {
	models.BudgetSectors
	models.BudgetActions
	models.Commissions
	models.Cities
	models.Communities
}

// GetSettings handles the get requests to give all datats in one batch
func GetSettings(ctx iris.Context) {
	var resp settingsResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.BudgetSectors.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration, secteurs budgétaires : " + err.Error()})
		return
	}
	if err := resp.BudgetActions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration, actions budgétaires : " + err.Error()})
		return
	}
	if err := resp.Commissions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration, commissions : " + err.Error()})
		return
	}
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration, villes : " + err.Error()})
		return
	}
	if err := resp.Communities.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration, interco : " + err.Error()})
		return
	}
	ctx.JSON(resp)
	ctx.StatusCode(http.StatusOK)
}
