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
