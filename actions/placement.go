package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// BatchPlacements handles the post request to insert ou update the database
// with a set of datas
func BatchPlacements(ctx iris.Context) {
	var req models.Placements
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de stages, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de stages, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de stages importé"})
}

// GetPlacements hadles the get request to fetch all placements from database
func GetPlacements(ctx iris.Context) {
	var resp models.Placements
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des stages, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetBeneficiaryPlacements hadles the get request to fetch all placements from database
func GetBeneficiaryPlacements(ctx iris.Context) {
	var resp models.Placements
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Stages d'un bénéficiaire, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetByBeneficiary(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Stages d'un bénéficiaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetBeneficiaryGroupPlacements hadles the get request to fetch all placements
// linked to beneficiaries that belong to a group from database
func GetBeneficiaryGroupPlacements(ctx iris.Context) {
	var resp models.Placements
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Stagiaires d'un groupe de bénéficiaires, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetByBeneficiaryGroup(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Stagiaires d'un groupe de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type updatePlacementReq struct {
	models.Placement `json:"Placement"`
}

// UpdatePlacement handles the put query to change the comment of a placement
func UpdatePlacement(ctx iris.Context) {
	var req updatePlacementReq
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Engagement de stagiaires, paramètre : " + err.Error()})
		return
	}
	if err = ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Engagement de stagiaire, décodage : " + err.Error()})
		return
	}
	req.ID = ID
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des stages, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}
