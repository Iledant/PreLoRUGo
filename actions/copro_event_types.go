package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// coproEventTypeResp embeddes an CoproEventType for json export
type coproEventTypeResp struct {
	C models.CoproEventType `json:"CoproEventType"`
}

// CreateCoproEventType hadles the post request to create a new CoproEventType into
// database
func CreateCoproEventType(ctx iris.Context) {
	var req coproEventTypeResp
	if err := ctx.ReadJSON(&req.C); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de type d'événement Copro, décodage : " + err.Error()})
		return
	}
	if err := req.C.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de type d'événement Copro : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.C.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de type d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateCoproEventType hadles the put request to modify an existing CoproEventType
func UpdateCoproEventType(ctx iris.Context) {
	var req coproEventTypeResp
	if err := ctx.ReadJSON(&req.C); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de type d'événement Copro, décodage : " + err.Error()})
		return
	}
	if err := req.C.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de type d'événement Copro : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.C.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de type d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteCoproEventType hadles the delete request to remove an existing CoproEventType
func DeleteCoproEventType(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de type d'événement Copro, paramètre : " + err.Error()})
		return
	}
	r := models.CoproEventType{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := r.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de type d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Type d'événement Copro supprimé"})
}

// GetCoproEventType hadles the get request to fetch an existing CoproEventType
// whose ID is given
func GetCoproEventType(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de type d'événement Copro, paramètre : " + err.Error()})
		return
	}
	var resp coproEventTypeResp
	resp.C.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.C.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de type d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetCoproEventTypes hadles the get request to fetch all CoproEventType from database
func GetCoproEventTypes(ctx iris.Context) {
	var resp models.CoproEventTypes
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération des types d'événement Copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
