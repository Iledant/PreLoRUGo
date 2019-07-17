package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// rpEventResp embeddes an RPEvent for json export
type rpEventResp struct {
	R models.RPEvent `json:"RPEvent"`
}

// CreateRPEvent hadles the post request to create a new RPEvent into
// database
func CreateRPEvent(ctx iris.Context) {
	var req rpEventResp
	if err := ctx.ReadJSON(&req.R); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'événement RP, décodage : " + err.Error()})
		return
	}
	if err := req.R.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'événement RP : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.R.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateRPEvent hadles the put request to modify an existing RPEvent
func UpdateRPEvent(ctx iris.Context) {
	var req rpEventResp
	if err := ctx.ReadJSON(&req.R); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'événement RP, décodage : " + err.Error()})
		return
	}
	if err := req.R.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification d'événement RP : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.R.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteRPEvent hadles the delete request to remove an existing RPEvent
func DeleteRPEvent(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'événement RP, paramètre : " + err.Error()})
		return
	}
	r := models.RPEvent{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := r.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Événement RP supprimé"})
}

// GetRPEvent hadles the get request to fetch an existing RPEvent
// whose ID is given
func GetRPEvent(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération d'événement RP, paramètre : " + err.Error()})
		return
	}
	var resp rpEventResp
	resp.R.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.R.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetRPEvents hadles the get request to fetch all RPEvent from database
func GetRPEvents(ctx iris.Context) {
	var resp models.RPEvents
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération d'événements RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
