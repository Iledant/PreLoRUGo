package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetPmtForecasts handles the get request to fetch payments ratios of a given year
func GetPmtForecasts(ctx iris.Context) {
	year, err := ctx.URLParamInt("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Prévisions de paiements, décodage : " + err.Error()})
		return
	}
	var resp models.PmtForecasts
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db, year); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Prévisions de paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
