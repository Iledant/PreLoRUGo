package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type reservationFeeReq struct {
	ReservationFee models.ReservationFee `json:"ReservationFee"`
}

// CreateReservationFee handles the post request to create a reservation fee
func CreateReservationFee(ctx iris.Context) {
	var req reservationFeeReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de réservation de logement, décodage : " +
			err.Error()})
		return
	}
	if err := req.ReservationFee.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de réservation de logement, paramètre : " +
			err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.ReservationFee.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de réservation de logement, requête : " +
			err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateReservationFee handles the put request to change a reservation fee
func UpdateReservationFee(ctx iris.Context) {
	var req reservationFeeReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de réservation de logement, décodage : " +
			err.Error()})
		return
	}
	if err := req.ReservationFee.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de réservation de logement, paramètre : " +
			err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.ReservationFee.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de réservation de logement, requête : " +
			err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// DeleteReservationFee handles the delete request to delete a housing transfer
func DeleteReservationFee(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression de réservation de logement, paramètre : " +
			err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	b := models.ReservationFee{ID: ID}
	if err := b.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de réservation de logement, requête : " +
			err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Réservation de logement supprimée"})
}

// BatchReservationFee handle the post request of a batch of reservation fees
func BatchReservationFee(ctx iris.Context) {
	var req models.ReservationFeeBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Batch de réservation de logement, décodage : " +
			err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	resp, err := req.Save(false, db)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de réservation de logement, requête : " +
			err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// TestBatchReservationFee handle the post request of a batch of reservation fees
func TestBatchReservationFee(ctx iris.Context) {
	var req models.ReservationFeeBatch
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Test batch de réservation de logement, décodage : " +
			err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	resp, err := req.Save(true, db)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Test batch de réservation de logement, requête : " +
			err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetPaginatedReservationFees handle the get request for reservations fees that
//  match a given search pattern returning a ReservationFee array, page number
//  and total count of items.
func GetPaginatedReservationFees(ctx iris.Context) {
	page, err := ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Page de réservation de logements, décodage Page : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.PaginatedQuery{Page: page, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.PaginatedReservationFees
	if err := resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de réservation de logements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type initialPaginatedReservationFeesResp struct {
	models.PaginatedReservationFees
	models.Cities
	models.Beneficiaries
	models.HousingTypologies
	models.HousingComments
	models.HousingTransfers
	models.ConventionTypes
	models.ReservationReports
}

// GetInitialPaginatedReservationFees handle the get request for reservations fees that
//  match a given search pattern returning a ReservationFee array, page number
//  and total count of items with all others datas needed for the frontend page.
func GetInitialPaginatedReservationFees(ctx iris.Context) {
	page, err := ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, décodage Page : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.PaginatedQuery{Page: page, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp initialPaginatedReservationFeesResp
	if err := resp.PaginatedReservationFees.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête fees : " + err.Error()})
		return
	}
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête cities : " + err.Error()})
		return
	}
	if err := resp.Beneficiaries.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête beneficiaries : " + err.Error()})
		return
	}
	if err := resp.HousingTypologies.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête typologies : " + err.Error()})
		return
	}
	if err := resp.HousingComments.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête comments : " + err.Error()})
		return
	}
	if err := resp.HousingTransfers.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête transfers : " + err.Error()})
		return
	}
	if err := resp.ConventionTypes.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête conventions types : " + err.Error()})
		return
	}
	if err := resp.ReservationReports.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page initiale de réservation de logements, requête report de réservations : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// ExportReservationFees handles the get request to fetch all reservation fees
// that matches the search pattern
func ExportReservationFees(ctx iris.Context) {
	search := ctx.URLParam("Search")
	req := models.PaginatedQuery{Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.ExportedReservationFees
	if err := resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Export de réservation de logements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

type reservationFeeSettingsResp struct {
	models.HousingComments
	models.HousingTransfers
	models.HousingTypologies
	models.ConventionTypes
}

// GetReservationFeeSettings handles the get request to fetch all datas for the
// reservation fee settings page
func GetReservationFeeSettings(ctx iris.Context) {
	db := ctx.Values().Get("db").(*sql.DB)
	var resp reservationFeeSettingsResp
	if err := resp.HousingComments.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration des réservations, requête commentaires : " + err.Error()})
		return
	}
	if err := resp.HousingTransfers.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration des réservations, requête cessions : " + err.Error()})
		return
	}
	if err := resp.HousingTypologies.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration des réservations, requête typologies : " + err.Error()})
		return
	}
	if err := resp.ConventionTypes.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Administration des réservations, requête types de conventions : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
