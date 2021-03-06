package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetAllPaymentDemands handle the get request to fetch all payment demands
// out of database
func GetAllPaymentDemands(ctx iris.Context) {
	var resp models.PaymentDemands
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Demandes de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type paymentDemandReq struct {
	PaymentDemand models.PaymentDemand `json:"PaymentDemand"`
}

// UpdatePaymentDemand handle the put request to update the excluded fields
func UpdatePaymentDemand(ctx iris.Context) {
	var req paymentDemandReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Mise à jour de demande de paiement, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.PaymentDemand.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Mise à jour de demande de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// BatchPaymentDemands handles the post request to update the database with a
// batch of payment demands
func BatchPaymentDemands(ctx iris.Context) {
	var req models.PaymentDemandBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Batch de demandes de paiement, décodage : " + err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Batch de demandes de paiement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de demandes de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de demande de paiement importé"})
}

// GetPaymentDemandCounts handle the get request to fetches the unprocessed
// or uncontrolled payment demands of the past 31 days
func GetPaymentDemandCounts(ctx iris.Context) {
	var resp models.PaymentDemandCounts
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Nombre de demandes de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetPaymentDemandStocks handle the get request to fetches the unprocessed
// or uncontrolled payment demands of the past 30 days
func GetPaymentDemandStocks(ctx iris.Context) {
	var resp models.PaymentDemandsStocks
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Stocks de demandes de paiement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
