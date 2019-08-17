package actions

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// renewProjectReq is used to embed an RenewProject in a post or put request
type renewProjectReq struct {
	RenewProject models.RenewProject `json:"RenewProject"`
}

// renewProjectDataResp embeddes the data for the get renew project datas request
type renewProjectDataResp struct {
	RenewProject models.RenewProject `json:"RenewProject"`
	models.RPLinkedCommitments
	models.Payments
	models.RenewProjectForecasts
	models.Commissions
	models.BudgetActions
	models.RPEventTypes
	models.FullRPEvents
	models.RPCmtCityJoins
	models.PreProgs
}

// GetRenewProjectDatas handles the get request to get renew project fields and
// related datas
func GetRenewProjectDatas(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("ID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, paramètre : " + err.Error()})
		return
	}
	var resp renewProjectDataResp
	db := ctx.Values().Get("db").(*sql.DB)
	resp.RenewProject.ID = ID
	if err = resp.RenewProject.GetByID(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête renewProject : " + err.Error()})
		return
	}
	if err = resp.RPLinkedCommitments.Get(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête commitments : " + err.Error()})
		return
	}
	if err = resp.Payments.GetLinkedToRenewProject(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête commitments : " + err.Error()})
		return
	}
	if err = resp.RenewProjectForecasts.Get(ID, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête forecasts : " + err.Error()})
		return
	}
	if err = resp.Commissions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête commissions : " + err.Error()})
		return
	}
	if err = resp.BudgetActions.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête actions : " + err.Error()})
		return
	}
	if err = resp.RPEventTypes.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête RPEventType : " + err.Error()})
		return
	}
	if err = resp.FullRPEvents.GetLinked(db, ID); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête FullRPEvents : " + err.Error()})
		return
	}
	if err = resp.RPCmtCityJoins.GetLinked(db, ID); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête lien engagements ville : " + err.Error()})
		return
	}
	year := (int64)(time.Now().Year())
	if err = resp.PreProgs.GetAllOfKind(year, models.KindRenewProject, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Datas de projet de renouvellement, requête lien engagements ville : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// CreateRenewProject handles the post request to create a renew project
func CreateRenewProject(ctx iris.Context) {
	var req renewProjectReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de projet de renouvellement, décodage : " + err.Error()})
		return
	}
	if err := req.RenewProject.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création de projet de renouvellement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RenewProject.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création de projet de renouvellement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// UpdateRenewProject handles the put request to modify a renew project
func UpdateRenewProject(ctx iris.Context) {
	var req renewProjectReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de projet de renouvellement, décodage : " + err.Error()})
		return
	}
	if err := req.RenewProject.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification de projet de renouvellement : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := req.RenewProject.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification de projet de renouvellement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(req)
}

// renewProjectsResp embeddes the renew projects and the cities to the frontend
// in a single request
type renewProjectsResp struct {
	models.Cities
	models.RenewProjects
	models.RPEventTypes
}

// GetRenewProjects handles the get request to handle all renew projets. In order
// to avoid multiples request from the front end, it also includes the list of
// the cities
func GetRenewProjects(ctx iris.Context) {
	var resp renewProjectsResp
	db := ctx.Values().Get("db").(*sql.DB)
	if err := resp.RenewProjects.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des projets de renouvellement, requête RU : " + err.Error()})
		return
	}
	if err := resp.Cities.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des projets de renouvellement, requête villes : " + err.Error()})
		return
	}
	if err := resp.RPEventTypes.GetAll(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Liste des projets de renouvellement, requête villes : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}

// DeleteRenewProject handles the delete request toremove a renew project from database
func DeleteRenewProject(ctx iris.Context) {
	ID, err := ctx.Params().GetInt64("rpID")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de projet de renouvellement, paramètre : " + err.Error()})
		return
	}
	rp := models.RenewProject{ID: ID}
	db := ctx.Values().Get("db").(*sql.DB)
	if err = rp.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression de projet de renouvellement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Projet de renouvellement supprimé"})
}

// BatchRenewProjects handle the post request to update and insert a batch of
// renew projects into the database
func BatchRenewProjects(ctx iris.Context) {
	var rp models.RenewProjectBatch
	if err := ctx.ReadJSON(&rp); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Batch de projets de renouvellement, décodage : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := rp.Save(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Batch de projets de renouvellement, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonMessage{"Batch de projets de renouvellement importé"})
}
