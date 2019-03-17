package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// BudgetSectorReq is used to embed a BudgetSector for request
type BudgetSectorReq struct {
	BudgetSector models.BudgetSector `json:"BudgetSector"`
}

// CreateBudgetSector handles the post request to create a new budget_sector
func CreateBudgetSector(ctx iris.Context) {
	var req BudgetSectorReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de secteur budgétaire, décodage : " + err.Error()})
		return
	}
	if err := req.BudgetSector.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de secteur budgétaire : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BudgetSector.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de secteur budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateBudgetSector handles the put request to modify a new budget_sector
func UpdateBudgetSector(ctx iris.Context) {
	var req BudgetSectorReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de secteur budgétaire, décodage : " + err.Error()})
		return
	}
	if err := req.BudgetSector.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de secteur budgétaire : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BudgetSector.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de secteur budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetBudgetSector handles the get request to fetch a budget_sector
func GetBudgetSector(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de secteur budgétaire, paramètre : " + err.Error()})
		return
	}
	var resp BudgetSectorReq
	resp.BudgetSector.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.BudgetSector.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Récupération de secteur budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetBudgetSectors handles the get request to fetch all budget_sectors
func GetBudgetSectors(ctx iris.Context) {
	var resp models.BudgetSectors
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des iDs, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteBudgetSector handles the get request to fetch all budget_sectors
func DeleteBudgetSector(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de secteur budgétaire, paramètre : " + err.Error()})
		return
	}
	resp := models.BudgetSector{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de secteur budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Logement supprimé"})
}
