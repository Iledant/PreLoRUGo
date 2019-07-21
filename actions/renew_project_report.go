package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"

	"github.com/kataras/iris"
)

// GetRenewProjectReport handle the get request to fetch the report on renew projets
func GetRenewProjectReport(ctx iris.Context) {
	var resp models.RenewProjectReport
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Report RU : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
