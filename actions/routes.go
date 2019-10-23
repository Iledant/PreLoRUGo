package actions

import (
	"database/sql"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

// SetRoutes initialize all routes for the application
func SetRoutes(app *iris.Application, superAdminEmail string, db *sql.DB) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	crsParty := app.Party("/api", crs).AllowMethods(iris.MethodOptions)

	crsParty.Post("/user/sign_up", setDBMiddleware(db, superAdminEmail), SignUp)
	crsParty.Post("/user/login", setDBMiddleware(db, superAdminEmail), Login)
	api := crsParty.Party("", setDBMiddleware(db, superAdminEmail))

	adminParty := api.Party("", AdminMiddleware)
	adminParty.Post("/user", CreateUser)
	adminParty.Put("/user/{userID}", UpdateUser)
	adminParty.Delete("/user/{userID}", DeleteUser)
	adminParty.Get("/users", GetUsers)

	adminParty.Post("/copro", CreateCopro)
	adminParty.Put("/copro", ModifyCopro)
	adminParty.Delete("/copro/{CoproID:int64}", DeleteCopro)
	adminParty.Post("/copros", BatchCopros)

	adminParty.Post("/budget_action", CreateBudgetAction)
	adminParty.Put("/budget_action", UpdateBudgetAction)
	adminParty.Delete("/budget_action/{baID}", DeleteBudgetAction)

	adminParty.Post("/renew_project", CreateRenewProject)
	adminParty.Put("/renew_project", UpdateRenewProject)
	adminParty.Delete("/renew_project/{rpID}", DeleteRenewProject)
	adminParty.Post("/renew_projects", BatchRenewProjects)

	adminParty.Post("/housing", CreateHousing)
	adminParty.Put("/housing", UpdateHousing)
	adminParty.Delete("/housing/{ID}", DeleteHousing)
	adminParty.Post("/housings", BatchHousings)

	adminParty.Post("/commitments", BatchCommitments)
	adminParty.Post("/commitments/link", LinkCommitment)
	adminParty.Post("/commitments/unlink", UnlinkCommitment)
	adminParty.Get("/commitments/forecasts", GetCmtForecasts)

	adminParty.Post("/payments", BatchPayments)
	adminParty.Get("/payments/forecasts", GetPmtForecasts)

	adminParty.Post("/budget_sector", CreateBudgetSector)
	adminParty.Put("/budget_sector", UpdateBudgetSector)
	adminParty.Delete("/budget_sector/{ID}", DeleteBudgetSector)

	adminParty.Post("/commission", CreateCommission)
	adminParty.Put("/commission", UpdateCommission)
	adminParty.Delete("/commission/{ID}", DeleteCommission)

	adminParty.Post("/community", CreateCommunity)
	adminParty.Put("/community", UpdateCommunity)
	adminParty.Delete("/community/{ID}", DeleteCommunity)
	adminParty.Post("/communities", BatchCommunities)

	adminParty.Post("/department", CreateDepartment)
	adminParty.Put("/department", UpdateDepartment)
	adminParty.Delete("/department/{ID}", DeleteDepartment)

	adminParty.Post("/city", CreateCity)
	adminParty.Put("/city", UpdateCity)
	adminParty.Delete("/city/{ID}", DeleteCity)
	adminParty.Post("/cities", BatchCities)

	adminParty.Post("/renew_project_forecasts", BatchRenewProjectForecasts)

	adminParty.Post("/housing_forecasts", BatchHousingForecasts)

	adminParty.Post("/copro_forecasts", BatchCoproForecasts)

	adminParty.Get("/settings", GetSettings)

	adminParty.Post("/ratios", BatchPmtRatios)
	adminParty.Get("/ratios/years", GetPmtRatiosYears)

	adminParty.Get("/pre_prog", GetPreProgs)

	adminParty.Post("/prog", SetProg)

	adminParty.Post("/rpls", CreateRPLS)
	adminParty.Put("/rpls", UpdateRPLS)
	adminParty.Delete("/rpls/{ID}", DeleteRPLS)
	adminParty.Post("/rpls/batch", BatchRPLS)
	adminParty.Get("/rpls/datas", GetRPLSDatas)

	adminParty.Post("/payment_credits", BatchPaymentCredits)

	coproUserParty := api.Party("", CoproMiddleware)
	coproUserParty.Post("/copro_forecast", CreateCoproForecast)
	coproUserParty.Put("/copro_forecast", UpdateCoproForecast)
	coproUserParty.Delete("/copro_forecast/{ID}", DeleteCoproForecast)
	coproUserParty.Post("/copro/commitments", LinkCommitmentsCopros)

	coproUserParty.Post("/copro_event_type", CreateCoproEventType)
	coproUserParty.Put("/copro_event_type", UpdateCoproEventType)
	coproUserParty.Delete("/copro_event_type/{ID}", DeleteCoproEventType)

	coproUserParty.Post("/copro_event", CreateCoproEvent)
	coproUserParty.Put("/copro_event", UpdateCoproEvent)
	coproUserParty.Delete("/copro_event/{ID}", DeleteCoproEvent)

	coproUserParty.Get("/pre_prog/copro", GetCoproPreProgs)
	coproUserParty.Post("/pre_prog/copro", SetCoproPreProgs)

	coproUserParty.Post("/copro/{CoproID}/copro_doc", CreateCoproDoc)
	coproUserParty.Put("/copro/{CoproID}/copro_doc", UpdateCoproDoc)
	coproUserParty.Delete("/copro/{CoproID}/copro_doc/{ID}", DeleteCoproDoc)

	renewProjectUserParty := api.Party("", RenewProjectMiddleware)
	renewProjectUserParty.Post("/renew_project_forecast", CreateRenewProjectForecast)
	renewProjectUserParty.Put("/renew_project_forecast", UpdateRenewProjectForecast)
	renewProjectUserParty.Delete("/renew_project_forecast/{ID}", DeleteRenewProjectForecast)

	renewProjectUserParty.Post("/rp_event_type", CreateRPEventType)
	renewProjectUserParty.Put("/rp_event_type", UpdateRPEventType)
	renewProjectUserParty.Delete("/rp_event_type/{ID}", DeleteRPEventType)

	renewProjectUserParty.Post("/rp_event", CreateRPEvent)
	renewProjectUserParty.Put("/rp_event", UpdateRPEvent)
	renewProjectUserParty.Delete("/rp_event/{ID}", DeleteRPEvent)

	renewProjectUserParty.Post("/rp_cmt_city_join", CreateRPCmtCityJoin)
	renewProjectUserParty.Put("/rp_cmt_city_join", UpdateRPCmtCityJoin)
	renewProjectUserParty.Delete("/rp_cmt_city_join/{ID}", DeleteRPCmtCityJoin)

	renewProjectUserParty.Get("/pre_prog/renew_project", GetRPPreProgs)
	renewProjectUserParty.Post("/pre_prog/renew_project", SetRPPreProgs)

	housingUserParty := api.Party("", HousingMiddleware)
	housingUserParty.Post("/housing_forecast", CreateHousingForecast)
	housingUserParty.Put("/housing_forecast", UpdateHousingForecast)
	housingUserParty.Delete("/housing_forecast/{ID}", DeleteHousingForecast)
	housingUserParty.Post("/housing/commitments", LinkCommitmentsHousings)

	housingUserParty.Get("/pre_prog/housing", GetHousingPreProgs)
	housingUserParty.Post("/pre_prog/housing", SetHousingPreProgs)

	housingUserParty.Post("/housing_summary", BatchHousingSummary)

	userParty := api.Party("", ActiveMiddleware)
	userParty.Post("/user/password", ChangeUserPwd)
	userParty.Post("/user/logout", Logout)
	userParty.Get("/budget_actions", GetBudgetActions)

	userParty.Get("/copro", GetCopros)
	userParty.Get("/copro/{ID}/datas", GetCoproDatas)

	userParty.Get("/renew_projects", GetRenewProjects)
	userParty.Get("/renew_project/{ID}/datas", GetRenewProjectDatas)

	userParty.Get("/housing/{ID}", GetHousing)
	userParty.Get("/housing/{ID}/datas", GetHousingDatas)
	userParty.Get("/housings", GetHousings)
	userParty.Get("/housings/datas", GetHousingsDatas)
	userParty.Get("/housings/paginated", GetPaginatedHousings)

	userParty.Get("/commitments", GetCommitments)
	userParty.Get("/commitments/paginated", GetPaginatedCommitments)
	userParty.Get("/commitments/unlinked", GetUnlinkedCommitments)
	userParty.Get("/commitments/export", ExportCommitments)

	userParty.Get("/beneficiaries", GetBeneficiaries)
	userParty.Get("/beneficiaries/paginated", GetPaginatedBeneficiaries)
	userParty.Get("/beneficiary/{ID}/datas", GetPaginatedBeneficiaryDatas)
	userParty.Get("/beneficiary/{ID}/export", GetExportBeneficiaryDatas)
	userParty.Get("/beneficiary/{ID}/payments", GetBeneficiaryPayments)

	userParty.Get("/payments", GetPayments)
	userParty.Get("/payments/paginated", GetPaginatedPayments)
	userParty.Get("/payments/export", GetExportedPayments)

	userParty.Get("/budget_sectors", GetBudgetSectors)
	userParty.Get("/budget_sector/{ID}", GetBudgetSector)

	userParty.Get("/community/{ID}", GetCommunity)
	userParty.Get("/communities", GetCommunities)

	userParty.Get("/department/{ID}", GetDepartment)
	userParty.Get("/departments", GetDepartments)

	userParty.Get("/commission/{ID}", GetCommission)
	userParty.Get("/commissions", GetCommissions)

	userParty.Get("/city/{ID}", GetCity)
	userParty.Get("/cities", GetCities)
	userParty.Get("/cities/paginated", GetPaginatedCities)

	userParty.Get("/renew_project_forecast/{ID}", GetRenewProjectForecast)
	userParty.Get("/renew_project_forecasts", GetRenewProjectForecasts)

	userParty.Get("/housing_forecast/{ID}", GetHousingForecast)
	userParty.Get("/housing_forecasts", GetHousingForecasts)

	userParty.Get("/copro_forecast/{ID}", GetCoproForecast)
	userParty.Get("/copro_forecasts", GetCoproForecasts)

	userParty.Get("/copro_event_types", GetCoproEventTypes)
	userParty.Get("/copro_event_type/{ID}", GetCoproEventType)

	userParty.Get("/copro_events", GetCoproEvents)
	userParty.Get("/copro_event/{ID}", GetCoproEvent)

	userParty.Get("/home", GetHome)

	userParty.Get("/ratios", GetPmtRatios)

	userParty.Get("/rp_event_types", GetRPEventTypes)
	userParty.Get("/rp_event_type/{ID}", GetRPEventType)

	userParty.Get("/rp_events", GetRPEvents)
	userParty.Get("/rp_event/{ID}", GetRPEvent)

	userParty.Get("/renew_project/report", GetRenewProjectReport)
	userParty.Get("/renew_project/report_per_community", GetRPPerCommunityReport)

	userParty.Get("/rp_cmt_city_joins", GetRPCmtCityJoins)
	userParty.Get("/rp_cmt_city_join/{ID}", GetRPCmtCityJoin)

	userParty.Get("/department_report", GetDptReport)

	userParty.Get("/city_report", GetCityReport)

	userParty.Get("/prog", GetProg)
	userParty.Get("/prog/datas", GetProgDatas)
	userParty.Get("/prog/years", GetProgYears)

	userParty.Get("/rpls", GetAllRPLS)
	userParty.Get("/rpls/report", RPLSReport)
	userParty.Get("/rpls/detailed_report", RPLSDetailedReport)

	userParty.Get("/summaries/datas", GetSummariesDatas)

	userParty.Get("/copro/{CoproID}/copro_docs", GetCoproDocs)

	userParty.Get("/copro/report", GetCoproReport)

	userParty.Get("/renew_project/multi_annual_report", GetRPMultiAnnualReport)

	userParty.Get("/payment_credits", GetAllPaymentCredits)
}

// setDBMiddleware return a middleware to add db to context values
func setDBMiddleware(db *sql.DB, superAdminEmail string) func(iris.Context) {
	return func(ctx iris.Context) {
		ctx.Values().Set("db", db)
		ctx.Values().Set("superAdminEmail", superAdminEmail)
		ctx.Next()
	}
}
