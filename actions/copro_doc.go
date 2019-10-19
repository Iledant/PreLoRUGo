package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"

	"github.com/kataras/iris"
)

// GetCoproDocs handle the get request to fetch all documents linked to a copro
func GetCoproDocs(ctx iris.Context) {
	CoproID, err := ctx.Params().GetInt64("CoproID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Documents d'une copro, erreur CoproID : " + err.Error()})
		return
	}
	var resp models.CoproDocs
	db := ctx.Values().Get("db").(*sql.DB)
	if err = resp.GetAll(CoproID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Documents d'une copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type coproDocResp struct {
	CoproDoc models.CoproDoc `json:"CoproDoc"`
}

// CreateCoproDoc handle the post request to create a document linked to a copro
func CreateCoproDoc(ctx iris.Context) {
	CoproID, err := ctx.Params().GetInt64("CoproID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'un document copro, paramètre : " + err.Error()})
		return
	}
	var resp coproDocResp
	if err = ctx.ReadJSON(&resp); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'un document copro, décodage : " + err.Error()})
		return
	}
	resp.CoproDoc.CoproID = CoproID
	db := ctx.Values().Get("db").(*sql.DB)
	if err = resp.CoproDoc.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'un document copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(resp)
}

// UpdateCoproDoc handle the put request to modify a document linked to a copro
func UpdateCoproDoc(ctx iris.Context) {
	var resp coproDocResp
	if err := ctx.ReadJSON(&resp); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification d'un document copro, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.CoproDoc.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'un document copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteCoproDoc handles the delete request to remove a document linked to a copro
func DeleteCoproDoc(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression d'un document copro, paramètre : " + err.Error()})
		return
	}
	coproDoc := models.CoproDoc{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = coproDoc.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'un document copro, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Document supprimé"})
}
