package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// RPCmtCityJoinReq is used to embed aRPCmtCityJoin for requests
type RPCmtCityJoinReq struct {
	RPCmtCityJoin models.RPCmtCityJoin `json:"RPCmtCityJoin"`
}

// CreateRPCmtCityJoin handles the post request to create a new community
func CreateRPCmtCityJoin(ctx iris.Context) {
	var req RPCmtCityJoinReq
	if err := ctx.ReadJSON(&req.RPCmtCityJoin); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de lien engagement ville, décodage : " + err.Error()})
		return
	}
	if err := req.RPCmtCityJoin.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de lien engagement ville : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RPCmtCityJoin.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateRPCmtCityJoin handles the put request to modify a new community
func UpdateRPCmtCityJoin(ctx iris.Context) {
	var req RPCmtCityJoinReq
	if err := ctx.ReadJSON(&req.RPCmtCityJoin); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de lien engagement ville, décodage : " + err.Error()})
		return
	}
	if err := req.RPCmtCityJoin.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de lien engagement ville : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RPCmtCityJoin.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetRPCmtCityJoin handles the get request to fetch a community
func GetRPCmtCityJoin(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de lien engagement ville, paramètre : " + err.Error()})
		return
	}
	var resp RPCmtCityJoinReq
	resp.RPCmtCityJoin.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPCmtCityJoin.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetRPCmtCityJoins handles the get request to fetch all communities
func GetRPCmtCityJoins(ctx iris.Context) {
	var resp models.RPCmtCityJoins
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des liens engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteRPCmtCityJoin handles the get request to fetch all communities
func DeleteRPCmtCityJoin(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de lien engagement ville, paramètre : " + err.Error()})
		return
	}
	resp := models.RPCmtCityJoin{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Lien engagement ville supprimé"})
}
