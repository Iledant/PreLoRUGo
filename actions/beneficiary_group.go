package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type beneficiaryGroupReq struct {
	BeneficiaryGroup models.BeneficiaryGroup `json:"BeneficiaryGroup"`
}

// CreateBeneficiaryGroup handles the post request to create a beneficiary group
func CreateBeneficiaryGroup(ctx iris.Context) {
	var req beneficiaryGroupReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de groupe de bénéficiaires, décodage : " + err.Error()})
		return
	}
	if err := req.BeneficiaryGroup.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de groupe de bénéficiaires, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BeneficiaryGroup.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de groupe de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateBeneficiaryGroup handles the put request to update a beneficiary group
func UpdateBeneficiaryGroup(ctx iris.Context) {
	var req beneficiaryGroupReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de groupe de bénéficiaires, décodage : " + err.Error()})
		return
	}
	if err := req.BeneficiaryGroup.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de groupe de bénéficiaires, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BeneficiaryGroup.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de groupe de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteBeneficiaryGroup handles the delete request to delete a beneficiary group
func DeleteBeneficiaryGroup(ctx iris.Context) {
	var req beneficiaryGroupReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de groupe de bénéficiaires, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.BeneficiaryGroup.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de groupe de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

type beneficiaryBelongReq struct {
	BeneficiaryIDs []int64 `json:"BeneficiaryIDs"`
}

// SetBeneficiaryGroup handle the post request to set all beneficiary that
// belongs to a beneficiary group
func SetBeneficiaryGroup(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de groupe de bénéficiaires, paramètre : " + err.Error()})
		return
	}
	var req beneficiaryBelongReq
	if err = ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Fixation de groupe de bénéficiaires, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	b := models.BeneficiaryGroup{ID: ID}
	if err := b.Set(req.BeneficiaryIDs, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Fixation de groupe de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetBeneficiaryGroups handles the get request to get all beneficiary groups
func GetBeneficiaryGroups(ctx iris.Context) {
	var resp models.BeneficiaryGroups
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Get(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des groupes de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetBeneficiaryGroupItems handles the get request to get all beneficiary datas
// from a beneficiary group
func GetBeneficiaryGroupItems(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Liste des bénéficiaires d'un groupe, décodage : " + err.Error()})
		return
	}
	var resp models.Beneficiaries
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GroupGet(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des bénéficiaires d'un groupe, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
