package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetPaginatedBeneficiaryDatas handle the get request for datas attached to a
// beneficiary that match a given search pattern returning a PageSize items,
// page number and total count of items.
func GetPaginatedBeneficiaryDatas(ctx iris.Context) {
	ID, err := ctx.Params().GetInt("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de données bénéficiaire, erreur ID : " + err.Error()})
		return
	}
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de données bénéficiaire, décodage Year : " + err.Error()})
		return
	}
	page, err := ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de données bénéficiaire, décodage Page : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.PaginatedQuery{Year: year, Page: page, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.PaginatedBeneficiaryDatas
	if err := resp.Get(db, &req, ID); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de données bénéficiaire, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
