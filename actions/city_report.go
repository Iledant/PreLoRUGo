package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetCityReport handle the get request to fetches commitments and payments par
// department
func GetCityReport(ctx iris.Context) {
	firstYear, err := ctx.URLParamInt64("firstYear")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport par commune, décodage firstYear : " + err.Error()})
		return
	}
	lastYear, err := ctx.URLParamInt64("lastYear")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport par commune, décodage lastYear : " + err.Error()})
		return
	}
	inseeCode, err := ctx.URLParamInt64("inseeCode")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport par commune, décodage inseeCode : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.CityReport
	if err = resp.GetAll(db, inseeCode, firstYear, lastYear); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Rapport par commune, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
