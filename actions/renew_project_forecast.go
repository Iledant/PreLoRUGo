package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// RenewProjectForecastReq is used to embed aRenewProjectForecast for requests
type RenewProjectForecastReq struct {
	RenewProjectForecast models.RenewProjectForecast `json:"RenewProjectForecast"`
}

// CreateRenewProjectForecast handles the post request to create a new RenewProjectForecast
func CreateRenewProjectForecast(ctx iris.Context) {
	var req RenewProjectForecastReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de prévision RU, décodage : " + err.Error()})
		return
	}
	if err := req.RenewProjectForecast.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de prévision RU : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RenewProjectForecast.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de prévision RU, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateRenewProjectForecast handles the put request to modify a new RenewProjectForecast
func UpdateRenewProjectForecast(ctx iris.Context) {
	var req RenewProjectForecastReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de prévision RU, décodage : " + err.Error()})
		return
	}
	if err := req.RenewProjectForecast.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de prévision RU : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RenewProjectForecast.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de prévision RU, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetRenewProjectForecast handles the get request to fetch a RenewProjectForecast
func GetRenewProjectForecast(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de prévision RU, paramètre : " + err.Error()})
		return
	}
	var resp RenewProjectForecastReq
	resp.RenewProjectForecast.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RenewProjectForecast.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de prévision RU, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetRenewProjectForecasts handles the get request to fetch all RenewProjectForecasts
func GetRenewProjectForecasts(ctx iris.Context) {
	var resp models.RenewProjectForecasts
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des prévision RUs, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteRenewProjectForecast handles the get request to fetch all RenewProjectForecasts
func DeleteRenewProjectForecast(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de prévision RU, paramètre : " + err.Error()})
		return
	}
	resp := models.RenewProjectForecast{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de prévision RU, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Prévision RU supprimé"})
}

// BatchRenewProjectForecasts handle the post request to update and insert a batch
// of RenewProjectForecasts into the database
func BatchRenewProjectForecasts(ctx iris.Context) {
	var b models.RenewProjectForecastBatch
	if err := ctx.ReadJSON(&b); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Prévision RUs, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Prévision RUs, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de Prévision RUs importé"})
}
