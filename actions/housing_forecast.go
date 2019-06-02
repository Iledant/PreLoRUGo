package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// HousingForecastReq is used to embed aHousingForecast for requests
type HousingForecastReq struct {
	HousingForecast models.HousingForecast `json:"HousingForecast"`
}

// CreateHousingForecast handles the post request to create a new HousingForecast
func CreateHousingForecast(ctx iris.Context) {
	var req HousingForecastReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de prévision logement, décodage : " + err.Error()})
		return
	}
	if err := req.HousingForecast.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de prévision logement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.HousingForecast.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de prévision logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateHousingForecast handles the put request to modify a new HousingForecast
func UpdateHousingForecast(ctx iris.Context) {
	var req HousingForecastReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de prévision logement, décodage : " + err.Error()})
		return
	}
	if err := req.HousingForecast.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de prévision logement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.HousingForecast.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de prévision logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetHousingForecast handles the get request to fetch a HousingForecast
func GetHousingForecast(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de prévision logement, paramètre : " + err.Error()})
		return
	}
	var resp HousingForecastReq
	resp.HousingForecast.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.HousingForecast.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de prévision logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetHousingForecasts handles the get request to fetch all HousingForecasts
func GetHousingForecasts(ctx iris.Context) {
	var resp models.HousingForecasts
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des prévision logements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteHousingForecast handles the get request to fetch all HousingForecasts
func DeleteHousingForecast(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de prévision logement, paramètre : " + err.Error()})
		return
	}
	resp := models.HousingForecast{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de prévision logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Prévision logement supprimé"})
}

// BatchHousingForecasts handle the post request to update and insert a batch
// of HousingForecasts into the database
func BatchHousingForecasts(ctx iris.Context) {
	var b models.HousingForecastBatch
	if err := ctx.ReadJSON(&b); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Prévision logements, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Prévision logements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de Prévision logements importé"})
}
