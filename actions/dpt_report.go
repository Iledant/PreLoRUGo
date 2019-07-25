package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetDptReport handle the get request to fetches commitments and payments par
// department
func GetDptReport(ctx iris.Context) {
	firstYear, err := ctx.URLParamInt64("firstYear")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport par département, décodage firstYear : " + err.Error()})
		return
	}
	lastYear, err := ctx.URLParamInt64("lastYear")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport par département, décodage lastYear : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.DptReport
	if err = resp.GetAll(db, firstYear, lastYear); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Rapport par département, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
