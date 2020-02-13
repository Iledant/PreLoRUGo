package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type beneficiaryReq struct {
	Beneficiary models.Beneficiary `json:"Beneficiary"`
}

// GetBeneficiaries handles the get request to fetch all beneficiaries
func GetBeneficiaries(ctx iris.Context) {
	var resp models.Beneficiaries
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetPaginatedBeneficiaries handles the get request to fetch all beneficiaries that
// match the given pattern and return a paginated struct with beneficiaries, page number
// and total page count
func GetPaginatedBeneficiaries(ctx iris.Context) {
	var req models.PaginatedQuery
	var err error
	req.Page, err = ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de bénéficiaires, décodage Page : " + err.Error()})
		return
	}
	req.Search = ctx.URLParam("Search")
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.PaginatedBeneficiaries
	if err = resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de bénéficiaires, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// CreateBeneficiary handles the post request to create a new beneficiary
func CreateBeneficiary(ctx iris.Context) {
	var req beneficiaryReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de bénéficiaire, décodage : " + err.Error()})
		return
	}
	if err := req.Beneficiary.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de bénéficiaire : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Beneficiary.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de bénéficiaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateBeneficiary handles the put request to modify a beneficiary
func UpdateBeneficiary(ctx iris.Context) {
	var req beneficiaryReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de bénéficiaire, décodage : " + err.Error()})
		return
	}
	if err := req.Beneficiary.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de bénéficiaire : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Beneficiary.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de bénéficiaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteBeneficiary handles the put request to modify a beneficiary
func DeleteBeneficiary(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de bénéficiaire, paramètre : " + err.Error()})
		return
	}
	b := models.Beneficiary{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de bénéficiaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Bénéficiaire supprimé"})
}
