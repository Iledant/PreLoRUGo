package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// rpEventTypeResp embeddes an RPEventType for json export
type rpEventTypeResp struct {
	R models.RPEventType `json:"RPEventType"`
}

// CreateRPEventType hadles the post request to create a new RPEventType into
// database
func CreateRPEventType(ctx iris.Context) {
	var req rpEventTypeResp
	if err := ctx.ReadJSON(&req.R); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de type d'événement RP, décodage : " + err.Error()})
		return
	}
	if err := req.R.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de type d'événement RP : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.R.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de type d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateRPEventType hadles the put request to modify an existing RPEventType
func UpdateRPEventType(ctx iris.Context) {
	var req rpEventTypeResp
	if err := ctx.ReadJSON(&req.R); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de type d'événement RP, décodage : " + err.Error()})
		return
	}
	if err := req.R.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de type d'événement RP : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.R.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de type d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteRPEventType hadles the delete request to remove an existing RPEventType
func DeleteRPEventType(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de type d'événement RP, paramètre : " + err.Error()})
		return
	}
	r := models.RPEventType{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := r.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de type d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Type d'événement RP supprimé"})
}

// GetRPEventType hadles the get request to fetch an existing RPEventType
// whose ID is given
func GetRPEventType(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de type d'événement RP, paramètre : " + err.Error()})
		return
	}
	var resp rpEventTypeResp
	resp.R.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.R.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de type d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetRPEventTypes hadles the get request to fetch all RPEventType from database
func GetRPEventTypes(ctx iris.Context) {
	var resp models.RPEventTypes
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération des types d'événement RP, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
