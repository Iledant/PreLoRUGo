package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// RPCmtCiyJoinReq is used to embed aRPCmtCiyJoin for requests
type RPCmtCiyJoinReq struct {
	RPCmtCiyJoin models.RPCmtCiyJoin `json:"RPCmtCiyJoin"`
}

// CreateRPCmtCiyJoin handles the post request to create a new community
func CreateRPCmtCiyJoin(ctx iris.Context) {
	var req RPCmtCiyJoinReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de lien engagement ville, décodage : " + err.Error()})
		return
	}
	if err := req.RPCmtCiyJoin.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de lien engagement ville : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RPCmtCiyJoin.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateRPCmtCiyJoin handles the put request to modify a new community
func UpdateRPCmtCiyJoin(ctx iris.Context) {
	var req RPCmtCiyJoinReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de lien engagement ville, décodage : " + err.Error()})
		return
	}
	if err := req.RPCmtCiyJoin.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de lien engagement ville : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RPCmtCiyJoin.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetRPCmtCiyJoin handles the get request to fetch a community
func GetRPCmtCiyJoin(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de lien engagement ville, paramètre : " + err.Error()})
		return
	}
	var resp RPCmtCiyJoinReq
	resp.RPCmtCiyJoin.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPCmtCiyJoin.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetRPCmtCiyJoins handles the get request to fetch all communities
func GetRPCmtCiyJoins(ctx iris.Context) {
	var resp models.RPCmtCiyJoins
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des liens engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteRPCmtCiyJoin handles the get request to fetch all communities
func DeleteRPCmtCiyJoin(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de lien engagement ville, paramètre : " + err.Error()})
		return
	}
	resp := models.RPCmtCiyJoin{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de lien engagement ville, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Lien engagement ville supprimé"})
}
