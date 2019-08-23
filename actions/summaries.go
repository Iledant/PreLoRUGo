package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// summariesREsp embeddes datas needded byt the summaries frontpage
type summariesResp struct {
	models.Cities
	models.RPLSYears
}

// GetSummariesDatas handles the get request to fetch datas of the summaries
// frontend page
func GetSummariesDatas(ctx iris.Context) {
	var resp summariesResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Données de synthèse, requête Cities : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	if err := resp.RPLSYears.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Données de synthèse, requête Cities : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.JSON(resp)
	ctx.StatusCode(http.StatusOK)
}
