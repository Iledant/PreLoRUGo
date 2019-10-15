package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// coproEventResp embeddes an CoproEvent for json export
type coproEventResp struct {
	C models.CoproEvent `json:"CoproEvent"`
}

// CreateCoproEvent hadles the post request to create a new CoproEvent into
// database
func CreateCoproEvent(ctx iris.Context) {
	var req coproEventResp
	if err := ctx.ReadJSON(&req.C); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'événement Copro, décodage : " + err.Error()})
		return
	}
	if err := req.C.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'événement Copro : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.C.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateCoproEvent hadles the put request to modify an existing CoproEvent
func UpdateCoproEvent(ctx iris.Context) {
	var req coproEventResp
	if err := ctx.ReadJSON(&req.C); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'événement Copro, décodage : " + err.Error()})
		return
	}
	if err := req.C.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification d'événement Copro : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.C.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteCoproEvent hadles the delete request to remove an existing CoproEvent
func DeleteCoproEvent(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'événement Copro, paramètre : " + err.Error()})
		return
	}
	r := models.CoproEvent{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := r.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Événement Copro supprimé"})
}

// GetCoproEvent hadles the get request to fetch an existing CoproEvent
// whose ID is given
func GetCoproEvent(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération d'événement Copro, paramètre : " + err.Error()})
		return
	}
	var resp coproEventResp
	resp.C.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.C.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetCoproEvents hadles the get request to fetch all CoproEvent from database
func GetCoproEvents(ctx iris.Context) {
	var resp models.CoproEvents
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération d'événements Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
