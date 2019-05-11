package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetPmtRatios handles the get request to fetch payments ratios of a given year
func GetPmtRatios(ctx iris.Context) {
	year, err := ctx.URLParamInt("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Ratios de paiements, décodage : " + err.Error()})
		return
	}
	var resp models.PmtRatios
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db, year); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Ratios de paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// BatchPmtRatios handles the post request to save a batch of ratios of a given
// year
func BatchPmtRatios(ctx iris.Context) {
	var req models.PmtRatioBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Batch de ratios de paiement, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de ratios de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de ratios de paiement traité"})
}

// GetPmtRatiosYears handles the get request to fetch payments ratios of a given year
func GetPmtRatiosYears(ctx iris.Context) {
	var resp models.PmtRatiosYears
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Années des ratios des paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
