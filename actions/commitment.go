package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// GetCommitments handles the get request to fetch all commitments
func GetCommitments(ctx iris.Context) {
	var resp models.Commitments
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des engagements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// GetPaginatedCommitments handles the get request to fetch all commitments that
// match the given pattern and return a paginated struct with commitments, page number
// and total page count
func GetPaginatedCommitments(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page d'engagements, décodage Year : " + err.Error()})
		return
	}
	page, err := ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page d'engagements, décodage Page : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.PaginatedQuery{Year: year, Page: page, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.PaginatedCommitments
	if err := resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page d'engagements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// ExportCommitments handles the get request to fetch all commitments that
// match the given pattern and return a list of commitments with full names
func ExportCommitments(ctx iris.Context) {
	year, err := ctx.URLParamInt64("Year")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Export d'engagements, décodage Year : " + err.Error()})
		return
	}
	search := ctx.URLParam("Search")
	req := models.ExportQuery{Year: year, Search: search}
	db := ctx.Values().Get("db").(*sql.DB)
	var resp models.ExportedCommitments
	if err := resp.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Export d'engagements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// BatchCommitments handle the post request to update and insert a batch of commitments into the database
func BatchCommitments(ctx iris.Context) {
	var b models.CommitmentBatch
	if err := ctx.ReadJSON(&b); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Engagements, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Engagements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de Engagements importé"})
}
