package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

type housingReq struct {
	Housing models.Housing `json:"Housing"`
}

// CreateHousing handles the post request to create a new housing
func CreateHousing(ctx iris.Context) {
	var req housingReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de logement, décodage : " + err.Error()})
		return
	}
	if err := req.Housing.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de logement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Housing.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateHousing handles the put request to modify a new housing
func UpdateHousing(ctx iris.Context) {
	var req housingReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de logement, décodage : " + err.Error()})
		return
	}
	if err := req.Housing.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de logement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.Housing.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(req)
}

// GetHousings handles the get request to fetch all housings
func GetHousings(ctx iris.Context) {
	resp := models.Housings{}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des logements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteHousing handles the get request to fetch all housings
func DeleteHousing(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de logement, paramètre : " + err.Error()})
		return
	}
	resp := models.Housing{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de logement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Logement supprimé"})
}

// BatchHousings handle the post request to update and insert a batch of housings into the database
func BatchHousings(ctx iris.Context) {
	var b models.HousingBatch
	if err := ctx.ReadJSON(&b); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Logements, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := b.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de Logements, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de Logements importé"})
}

// PaginatedTedHousingsResp embeddes response for PaginatedHousing in order
// to optionally add Cities to response
type PaginatedTedHousingsResp struct {
	models.PaginatedHousings
	models.Cities
}

// GetPaginatedHousings handles the get request to fetch all housings that
// match the given pattern and return a paginated struct with housings, page
// number and total page count. An optional CitiesList flag is added to embed
// the datats in a single request/response
func GetPaginatedHousings(ctx iris.Context) {
	var req models.PaginatedQuery
	var err error
	req.Page, err = ctx.URLParamInt64("Page")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de logements, décodage Page : " + err.Error()})
		return
	}
	req.Search = ctx.URLParam("Search")
	db := ctx.Values().Get("db").(*sql.DB)
	var resp PaginatedTedHousingsResp
	if err = resp.PaginatedHousings.Get(db, &req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Page de logements, requête : " + err.Error()})
		return
	}
	cl, err := ctx.URLParamBool("CitiesList")
	if err == nil && cl {
		if err = resp.Cities.GetAll(db); err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(jsonError{"Page de logements, requête liste de villes : " + err.Error()})
			return
		}
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
