package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// CoproForecastReq is used to embed aCoproForecast for requests
type CoproForecastReq struct {
	CoproForecast models.CoproForecast `json:"CoproForecast"`
}

// CreateCoproForecast handles the post request to create a new CoproForecast
func CreateCoproForecast(ctx iris.Context) {
	var req CoproForecastReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de prévision copro, décodage : " + err.Error()})
		return
	}
	if err := req.CoproForecast.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de prévision copro : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.CoproForecast.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de prévision copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateCoproForecast handles the put request to modify a new CoproForecast
func UpdateCoproForecast(ctx iris.Context) {
	var req CoproForecastReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de prévision copro, décodage : " + err.Error()})
		return
	}
	if err := req.CoproForecast.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de prévision copro : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.CoproForecast.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de prévision copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetCoproForecast handles the get request to fetch a CoproForecast
func GetCoproForecast(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de prévision copro, paramètre : " + err.Error()})
		return
	}
	var resp CoproForecastReq
	resp.CoproForecast.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.CoproForecast.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de prévision copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetCoproForecasts handles the get request to fetch all CoproForecasts
func GetCoproForecasts(ctx iris.Context) {
	var resp models.CoproForecasts
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des prévision copros, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteCoproForecast handles the get request to fetch all CoproForecasts
func DeleteCoproForecast(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de prévision copro, paramètre : " + err.Error()})
		return
	}
	resp := models.CoproForecast{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de prévision copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Prévision copro supprimé"})
}

// BatchCoproForecasts handle the post request to update and insert a batch
// of CoproForecasts into the database
func BatchCoproForecasts(ctx iris.Context) {
	var b models.CoproForecastBatch
	if err := ctx.ReadJSON(&b); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Prévision copros, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Prévision copros, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de Prévision copros importé"})
}
