package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type conventionTypeReq struct {
	ConventionType models.ConventionType `json:"ConventionType"`
}

// CreateConventionType handles the post request to create a convention type
func CreateConventionType(ctx iris.Context) {
	var req conventionTypeReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de type de convention, décodage : " + err.Error()})
		return
	}
	if err := req.ConventionType.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de type de convention, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.ConventionType.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de type de convention, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateConventionType handles the put request to update a convention type
func UpdateConventionType(ctx iris.Context) {
	var req conventionTypeReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de type de convention, décodage : " + err.Error()})
		return
	}
	if err := req.ConventionType.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de type de convention, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.ConventionType.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de type de convention, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteConventionType handles the delete request to delete a convention type
func DeleteConventionType(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de type de convention, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	b := models.ConventionType{ID: ID}
	if err := b.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de type de convention, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Type de convention supprimé"})
}

// GetConventionTypes handles the get request to get all housing comment
func GetConventionTypes(ctx iris.Context) {
	var resp models.ConventionTypes
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des types de convention, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
