package actions

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// rplsDatasResp is used to sent all needed datas to frontend page
type rplsDatasResp struct {
	models.RPLSArray
	models.Cities
}

// GetRPLSDatas handle the get request to fetch all datas for dedicated frontend
// page
func GetRPLSDatas(ctx iris.Context) {
	var resp rplsDatasResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPLSArray.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Datas RPLS, requête RPLS : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Datas RPLS, requête Cities : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetAllRPLS handle the get request to fetch all RPLS
func GetAllRPLS(ctx iris.Context) {
	var resp models.RPLSArray
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Liste RPLS, requête : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// rplsResp is used to embed JSON response
type rplsResp struct {
	RPLS models.RPLS `json:"RPLS"`
}

// CreateRPLS handle the post request to create a new RPLS
func CreateRPLS(ctx iris.Context) {
	var resp rplsResp
	if err := ctx.ReadJSON(&resp.RPLS); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de RPLS, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPLS.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(resp)
}

// UpdateRPLS handle the put request to modify an existing RPLS
func UpdateRPLS(ctx iris.Context) {
	var resp rplsResp
	if err := ctx.ReadJSON(&resp.RPLS); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de RPLS, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RPLS.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteRPLS handles the delete request to remove an existing RPLS
func DeleteRPLS(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de RPLS, ID introuvable : " + err.Error()})
		return
	}
	req := models.RPLS{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = req.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{`RPLS supprimé`})
}

// BatchRPLS handle the post request to insert a batch of RPLS into database
func BatchRPLS(ctx iris.Context) {
	var req models.RPLSBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch RPLS, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch RPLS importé"})
}

func decodeParams(p *models.RPLSReportParams, ctx iris.Context) (err error) {
	if p.RPLSYear, err = ctx.URLParamInt64("RPLSYear"); err != nil {
		return fmt.Errorf("RPLSYear : %v", err)
	}
	if p.FirstYear, err = ctx.URLParamInt64("FirstYear"); err != nil {
		return fmt.Errorf("FirstYear : %v", err)
	}
	if p.LastYear, err = ctx.URLParamInt64("LastYear"); err != nil {
		return fmt.Errorf("LastYear : %v", err)
	}
	if p.RPLSMin, err = ctx.URLParamFloat64("RPLSMin"); err != nil {
		return fmt.Errorf("RPLSMin : %v", err)
	}
	if p.RPLSMax, err = ctx.URLParamFloat64("RPLSMax"); err != nil {
		return fmt.Errorf("RPLSMax : %v", err)
	}
	return nil
}

// RPLSReport handle the get request to fetch the RPLS Report
func RPLSReport(ctx iris.Context) {
	var req models.RPLSReportParams
	if err := decodeParams(&req, ctx); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport RPLS, décodage " + err.Error()})
		return
	}
	var resp models.RPLSReport
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(&req, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Rapport RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// RPLSDetailedReport handle the get request to fetch the RPLS Report
func RPLSDetailedReport(ctx iris.Context) {
	var req models.RPLSReportParams
	if err := decodeParams(&req, ctx); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Rapport détaillé RPLS, décodage " + err.Error()})
		return
	}
	var resp models.RPLSDetailedReport
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(&req, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Rapport détaillé RPLS, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
