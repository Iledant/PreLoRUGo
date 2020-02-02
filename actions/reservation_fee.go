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
		ctx.JSON(jsonError{"Création de réservation de logement, décodage : " + err.Error()})
		return
	}
	if err := req.ReservationFee.Valid(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de réservation de logement, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.ReservationFee.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de réservation de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}
