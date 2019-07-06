package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type getCoprosResp struct {
	models.Copros
	models.Cities
}

// GetCopros handles the get request to fetch all copros
func GetCopros(ctx iris.Context) {
	var resp getCoprosResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Copros.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des copropriétés, copros : " + err.Error()})
		return
	}
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des copropriétés, copros : " + err.Error()})
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
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de copropriété, décodage : " + err.Error()})
		return
	}
	if err := c.Copro.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de copropriété : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := c.Copro.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de copropriété, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(c)
}

// coproDatasResp embeddes the different datas for the get copro datas request
type coproDatasResp struct {
	Copro models.Copro `json:"Copro"`
	models.CoproLinkedCommitments
	models.Payments
	models.Commissions
	models.CoproForecasts
}

// GetCoproDatas handle the get request to fetch copro fields, commitments and
// payments linked to that copro
func GetCoproDatas(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données d'une copropriété, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp coproDatasResp
	resp.Copro.ID = ID
	if err = resp.Copro.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données d'une copropriété, requête copro : " + err.Error()})
		return
	}
	if err = resp.CoproLinkedCommitments.Get(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données d'une copropriété, requête commitment : " + err.Error()})
		return
	}
	if err = resp.Payments.GetLinkedToCopro(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données d'une copropriété, requête payment : " + err.Error()})
		return
	}
	if err = resp.Commissions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données d'une copropriété, requête commissions : " + err.Error()})
		return
	}
	if err = resp.CoproForecasts.Get(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données d'une copropriété, requête forecasts : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// ModifyCopro handles the put request to modify a copro
func ModifyCopro(ctx iris.Context) {
	var c coproReq
	if err := ctx.ReadJSON(&c); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de copropriété, décodage : " + err.Error()})
		return
	}
	if err := c.Copro.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de copropriété : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := c.Copro.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de copropriété, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(c)
}

// DeleteCopro handles the delete request to remove a copro from database
func DeleteCopro(ctx iris.Context) {
	var c models.Copro
	var err error
	if c.ID, err = ctx.Params().GetInt64("CoproID"); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de copropriété, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = c.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de copropriété, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Copropriété supprimée"})
}

// BatchCopros handle the post request to update and insert a batch of copros into the database
func BatchCopros(ctx iris.Context) {
	var c models.CoproBatch
	if err := ctx.ReadJSON(&c); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de copropriétés, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := c.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de copropriétés, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de copropriétés importé"})
}
