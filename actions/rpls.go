package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// rplsDatasResp is used to sent all needed datas to frontend page
type rplsDatasResp struct {
	models.RPLSArray
	models.Cities
}

// GetRPLSDatas handle the get request to fetch all datas for dedicated frontend
// page
func GetRPLSDatas(ctx iris.Context) {
	var resp rplsDatasResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPLSArray.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Datas RPLS, requête RPLS : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Datas RPLS, requête Cities : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetAllRPLS handle the get request to fetch all RPLS
func GetAllRPLS(ctx iris.Context) {
	var resp models.RPLSArray
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Liste RPLS, requête : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// rplsResp is used to embed JSON response
type rplsResp struct {
	RPLS models.RPLS `json:"RPLS"`
}

// CreateRPLS handle the post request to create a new RPLS
func CreateRPLS(ctx iris.Context) {
	var resp rplsResp
	if err := ctx.ReadJSON(&resp.RPLS); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de RPLS, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPLS.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(resp)
}

// UpdateRPLS handle the put request to modify an existing RPLS
func UpdateRPLS(ctx iris.Context) {
	var resp rplsResp
	if err := ctx.ReadJSON(&resp.RPLS); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de RPLS, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPLS.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteRPLS handles the delete request to remove an existing RPLS
func DeleteRPLS(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de RPLS, ID introuvable : " + err.Error()})
		return
	}
	req := models.RPLS{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = req.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{`RPLS supprimé`})
}

// BatchRPLS handle the post request to insert a batch of RPLS into database
func BatchRPLS(ctx iris.Context) {
	var req models.RPLSBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch RPLS, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch RPLS importé"})
}
