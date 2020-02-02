package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type housingConventionReq struct {
	HousingConvention models.HousingConvention `json:"HousingConvention"`
}

// CreateHousingConvention handles the post request to create a housing convention
func CreateHousingConvention(ctx iris.Context) {
	var req housingConventionReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de convention de logement, décodage : " + err.Error()})
		return
	}
	if err := req.HousingConvention.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de convention de logement, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.HousingConvention.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de convention de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateHousingConvention handles the put request to update a housing convention
func UpdateHousingConvention(ctx iris.Context) {
	var req housingConventionReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de convention de logement, décodage : " + err.Error()})
		return
	}
	if err := req.HousingConvention.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de convention de logement, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.HousingConvention.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de convention de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteHousingConvention handles the delete request to delete a housing convention
func DeleteHousingConvention(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de convention de logement, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	b := models.HousingConvention{ID: ID}
	if err := b.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de convention de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Convention de logement supprimée"})
}

// GetHousingConventions handles the get request to get all housing conventions
func GetHousingConventions(ctx iris.Context) {
	var resp models.HousingConventions
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des conventions de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
