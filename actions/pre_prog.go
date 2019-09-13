package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetPreProgs handles the get request to fetch all pre progs datas of a given
// year
func GetPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Préprogrammation d'une année, décodage : " + err.Error()})
		return
	}
	var resp models.PreProgs
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Préprogrammation d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetCoproPreProgs handles the get request to fetch the pre prog of copro
// operations of a given year
func GetCoproPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Préprogrammation copro d'une année, décodage : " + err.Error()})
		return
	}
	var resp models.FcPreProgs
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAllOfKind(year, models.KindCopro, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Préprogrammation copro d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetRPPreProgs handles the get request to fetch the pre prog of renew project
// operations of a given year
func GetRPPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Préprogrammation RU d'une année, décodage : " + err.Error()})
		return
	}
	var resp models.FcPreProgs
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAllOfKind(year, models.KindRenewProject, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Préprogrammation RU d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetHousingPreProgs handles the get request to fetch the pre prog of renew project
// operations of a given year
func GetHousingPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Préprogrammation logement d'une année, décodage : " + err.Error()})
		return
	}
	var resp models.FcPreProgs
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAllOfKind(year, models.KindHousing, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Préprogrammation logement d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// SetCoproPreProgs handles the post request to set the pre programmation of
// copro operation of a given year
func SetCoproPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la préprogrammation copro d'une année, décodage année : " + err.Error()})
		return
	}
	var req models.PreProgBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la préprogrammation copro d'une année, décodage batch : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(models.KindCopro, year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Fixation de la préprogrammation copro d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch importé"})
}

// SetRPPreProgs handles the post request to set the pre programmation of
// RP operation of a given year
func SetRPPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la préprogrammation RU d'une année, décodage année : " + err.Error()})
		return
	}
	var req models.PreProgBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la préprogrammation RU d'une année, décodage batch : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(models.KindRenewProject, year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Fixation de la préprogrammation RU d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch importé"})
}

// SetHousingPreProgs handles the post request to set the pre programmation of
// housing operation of a given year
func SetHousingPreProgs(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la préprogrammation logement d'une année, décodage année : " + err.Error()})
		return
	}
	var req models.PreProgBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de la préprogrammation logement d'une année, décodage batch : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(models.KindHousing, year, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Fixation de la préprogrammation logement d'une année, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch importé"})
}
