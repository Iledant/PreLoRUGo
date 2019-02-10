package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetCopros handles the get request to fetch all copros
func GetCopros(ctx iris.Context) {
	var resp models.Copros
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Liste des copropriétés : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// coproReq embeddes a copro model to handle create or modify request
type coproReq struct {
	Copro models.Copro `json:"Copro"`
}

// CreateCopro handles the post request to create a new copro
func CreateCopro(ctx iris.Context) {
	var c coproReq
	if err := ctx.ReadJSON(&c); err != nil {
		ctx.JSON(jsonError{"Création de copropriété, décodage : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	if err := c.Copro.Validate(); err != nil {
		ctx.JSON(jsonError{"Création de copropriété : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := c.Copro.Create(db); err != nil {
		ctx.JSON(jsonError{"Création de copropriété, requête : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.JSON(c)
	ctx.StatusCode(http.StatusOK)
}

// ModifyCopro handles the put request to modify a copro
func ModifyCopro(ctx iris.Context) {
	var c coproReq
	if err := ctx.ReadJSON(&c); err != nil {
		ctx.JSON(jsonError{"Modification de copropriété, décodage : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	if err := c.Copro.Validate(); err != nil {
		ctx.JSON(jsonError{"Modification de copropriété : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := c.Copro.Update(db); err != nil {
		ctx.JSON(jsonError{"Modification de copropriété, requête : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.JSON(c)
	ctx.StatusCode(http.StatusOK)
}

// DeleteCopro handles the delete request to remove a copro from database
func DeleteCopro(ctx iris.Context) {
	var c models.Copro
	var err error
	if c.ID, err = ctx.Params().GetInt64("CoproID"); err != nil {
		ctx.JSON(jsonError{"Suppression de copropriété, paramètre : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = c.Delete(db); err != nil {
		ctx.JSON(jsonError{"Modification de copropriété, requête : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.JSON(jsonMessage{"Copropriété supprimée"})
	ctx.StatusCode(http.StatusOK)
}

// BatchCopros handle the post request to update and insert a batch of copros into the database
func BatchCopros(ctx iris.Context) {
	var c models.CoproBatch
	if err := ctx.ReadJSON(&c); err != nil {
		ctx.JSON(jsonError{"Batch de copropriétés, décodage : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := c.Save(db); err != nil {
		ctx.JSON(jsonError{"Batch de copropriétés, requête : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.JSON(jsonMessage{"Batch de copropriétés importé"})
	ctx.StatusCode(http.StatusOK)
}
