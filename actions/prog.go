package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetProg handles the get request to fetch all programming datas of a given
// year
func GetProg(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Programmation d'une année, décodage : " + err.Error()})
		return
	}
	var resp models.Progs
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Programmation d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type progDatasResp struct {
	models.Progs
	models.Commissions
	models.BudgetActions
	models.Copros
	models.RenewProjects
}

// GetProgDatas handles the get request to fetch all datas for the frontend page
// dedicated to the programmation
func GetProgDatas(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Données de programmation d'une année, décodage : " + err.Error()})
		return
	}
	var resp progDatasResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Progs.GetAll(year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données de programmation d'une année, requête : " + err.Error()})
		return
	}
	if err := resp.Commissions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données de programmation d'une année, requête commissions : " + err.Error()})
		return
	}
	if err := resp.BudgetActions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données de programmation d'une année, requête actions budgétaires : " + err.Error()})
		return
	}
	if err := resp.Copros.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données de programmation d'une année, requête copros : " + err.Error()})
		return
	}
	if err := resp.RenewProjects.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Données de programmation d'une année, requête projets RU : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// SetProg handles the post request to set the programmation of a given year
func SetProg(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la programmation d'une année, décodage année : " + err.Error()})
		return
	}
	var req models.ProgBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la programmation d'une année, décodage batch : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Fixation de la programmation d'une année, requête : " + err.Error()})
		return
	}
	var resp models.Progs
	if err := resp.GetAll(year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Fixation de la programmation d'une année, requête programmation : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetProgYears handles the get request to fetch all programmation years from the
// database
func GetProgYears(ctx iris.Context) {
	var resp models.ProgYears
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Années de programmation, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
