package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetBudgetActions handles the get request to fetch all budget actions
func GetBudgetActions(ctx iris.Context) {
	var resp models.BudgetActions
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des actions budgétaires : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type budgetActionReq struct {
	BudgetAction models.BudgetAction `json:"BudgetAction"`
}

// CreateBudgetAction handles the post request to create a new budget action
func CreateBudgetAction(ctx iris.Context) {
	var req budgetActionReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'action budgétaire, décodage : " + err.Error()})
		return
	}
	if err := req.BudgetAction.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'action budgétaire : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BudgetAction.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'action budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateBudgetAction handles the post request to create a new budget action
func UpdateBudgetAction(ctx iris.Context) {
	var req budgetActionReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'action budgétaire, décodage : " + err.Error()})
		return
	}
	if err := req.BudgetAction.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification d'action budgétaire : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BudgetAction.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'action budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteBudgetAction handles the delete request to remove a budget action from database
func DeleteBudgetAction(ctx iris.Context) {
	baID, err := ctx.Params().GetInt64("baID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'action budgétaire, paramètre : " + err.Error()})
		return
	}
	budgetAction := models.BudgetAction{ID: baID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = budgetAction.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'action budgétaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Action budgétaire supprimée"})
}
