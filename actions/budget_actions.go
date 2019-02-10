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
		ctx.JSON(jsonError{"Liste des actions budgétaires : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
