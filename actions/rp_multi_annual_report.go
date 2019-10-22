package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetRPMultiAnnualReport handles the get request to fetch multiannual report
// from database
func GetRPMultiAnnualReport(ctx iris.Context) {
	var resp models.RPMultiAnnualReports
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"rapport sur les copropriétés, requête : " + err.Error()})
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
