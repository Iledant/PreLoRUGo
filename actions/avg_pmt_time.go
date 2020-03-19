package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetAvgPmtTimes handle the get request to fetches all average payment times
// for the past twelve monthes
func GetAvgPmtTimes(ctx iris.Context) {
	var resp models.AvgPmtTimes
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Durée moyenne de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)

}
