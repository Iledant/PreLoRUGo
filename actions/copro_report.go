package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetCoproReport handles the get request to fetch the copro report from database
func GetCoproReport(ctx iris.Context) {
	var resp models.CoproReports
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"rapport sur les copropriétés, requête : " + err.Error()})
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
