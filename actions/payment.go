package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type paymentReq struct {
	Payment models.Payment `json:"Payment"`
}

// GetPayments handles the get request to fetch all payments
func GetPayments(ctx iris.Context) {
	var resp models.Payments
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// BatchPayments handle the post request to update and insert a batch of payments into the database
func BatchPayments(ctx iris.Context) {
	var b models.PaymentBatch
	if err := ctx.ReadJSON(&b); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Paiements, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de Paiements importé"})
}

// GetPaginatedPayments handle the get request for commitments that match a given
// search pattern returning a PageSize items, page number and total count of items.
func GetPaginatedPayments(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de paiements, décodage Year : " + err.Error()})
		return
	}
	page, err := ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de paiements, décodage Page : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.PaginatedQuery{Year: year, Page: page, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.PaginatedPayments
	if err := resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetExportedPayments handle the get request for commitments that match a given
// search pattern returning a payments with full linked names.
func GetExportedPayments(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Export de paiements, décodage Year : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.ExportQuery{Year: year, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.ExportPayments
	if err := resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Export de paiements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
