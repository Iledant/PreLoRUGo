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

	adminParty := api.Party("", RightsMiddleWare(&admHandler))
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

	adminParty.Get("/pre_prog", GetPreProgs)

	adminParty.Post("/prog", SetProg)

	adminParty.Post("/rpls", CreateRPLS)
	adminParty.Put("/rpls", UpdateRPLS)
	adminParty.Delete("/rpls/{ID}", DeleteRPLS)
	adminParty.Post("/rpls/batch", BatchRPLS)
	adminParty.Get("/rpls/datas", GetRPLSDatas)

	adminParty.Post("/payment_credits", BatchPaymentCredits)

	adminParty.Post("/payment_credit_journal", BatchPaymentCreditJournals)

	adminParty.Post("/home_message", SetHomeMessage)

	adminParty.Post("/placements", BatchPlacements)
	adminParty.Put("/placement/{ID}", UpdatePlacement)

	adminParty.Post("/beneficiary_group", CreateBeneficiaryGroup)
	adminParty.Put("/beneficiary_group", UpdateBeneficiaryGroup)
	adminParty.Delete("/beneficiary_group/{ID}", DeleteBeneficiaryGroup)
	adminParty.Post("/beneficiary_group/{ID}", SetBeneficiaryGroup)

	adminParty.Post("/housing_typology", CreateHousingTypology)
	adminParty.Put("/housing_typology", UpdateHousingTypology)
	adminParty.Delete("/housing_typology/{ID}", DeleteHousingTypology)

	adminParty.Post("/housing_convention", CreateHousingConvention)
	adminParty.Put("/housing_convention", UpdateHousingConvention)
	adminParty.Delete("/housing_convention/{ID}", DeleteHousingConvention)

	adminParty.Post("/housing_comment", CreateHousingComment)
	adminParty.Put("/housing_comment", UpdateHousingComment)
	adminParty.Delete("/housing_comment/{ID}", DeleteHousingComment)

	adminParty.Post("/housing_transfer", CreateHousingTransfer)
	adminParty.Put("/housing_transfer", UpdateHousingTransfer)
	adminParty.Delete("/housing_transfer/{ID}", DeleteHousingTransfer)

	adminParty.Post("/convention_type", CreateConventionType)
	adminParty.Put("/convention_type", UpdateConventionType)
	adminParty.Delete("/convention_type/{ID}", DeleteConventionType)

	adminParty.Post("/beneficiary", CreateBeneficiary)
	adminParty.Put("/beneficiary", UpdateBeneficiary)
	adminParty.Delete("/beneficiary/{ID}", DeleteBeneficiary)

	adminParty.Post("/housing_type", CreateHousingType)
	adminParty.Put("/housing_type", UpdateHousingType)
	adminParty.Delete("/housing_type/{ID}", DeleteHousingType)

	adminParty.Post("/iris_housing_type", BatchIRISHousingType)

	adminParty.Get("/commitments/eldest", GetEldestCommitments)
	adminParty.Get("/commitments/unpaid", GetUnpaidCommitments)

	coproUserParty := api.Party("", RightsMiddleWare(&coproHandler))
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

	coproUserParty.Post("/copro/{CoproID}/copro_doc", CreateCoproDoc)
	coproUserParty.Put("/copro/{CoproID}/copro_doc", UpdateCoproDoc)
	coproUserParty.Delete("/copro/{CoproID}/copro_doc/{ID}", DeleteCoproDoc)

	coproPreProgParty := api.Party("", RightsMiddleWare(&coproPreProgHandler))
	coproPreProgParty.Post("/pre_prog/copro", SetCoproPreProgs)

	renewProjectUserParty := api.Party("", RightsMiddleWare(&rpHandler))
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

	renewProjectPreProgUserParty := api.Party("", RightsMiddleWare(&rpPreProgHandler))
	renewProjectPreProgUserParty.Post("/pre_prog/renew_project", SetRPPreProgs)

	housingUserParty := api.Party("", RightsMiddleWare(&housingHandler))
	housingUserParty.Post("/housing_forecast", CreateHousingForecast)
	housingUserParty.Put("/housing_forecast", UpdateHousingForecast)
	housingUserParty.Delete("/housing_forecast/{ID}", DeleteHousingForecast)
	housingUserParty.Post("/housing/commitments", LinkCommitmentsHousings)

	housingUserParty.Get("/pre_prog/housing", GetHousingPreProgs)

	housingUserParty.Post("/housing_summary", BatchHousingSummary)

	housingPreProgUserParty := api.Party("", RightsMiddleWare(&housingPreProgHandler))
	housingPreProgUserParty.Post("/pre_prog/housing", SetHousingPreProgs)

	reservationUserParty := api.Party("", RightsMiddleWare(&reservationHandler))
	reservationUserParty.Post("/reservation_fee", CreateReservationFee)
	reservationUserParty.Get("/reservation_fees", GetPaginatedReservationFees)
	reservationUserParty.Get("/reservation_fees/initial", GetInitialPaginatedReservationFees)
	reservationUserParty.Get("/reservation_fees/export", ExportReservationFees)
	reservationUserParty.Post("/reservation_fee/batch", BatchReservationFee)
	reservationUserParty.Post("/reservation_fee/batch/test", TestBatchReservationFee)
	reservationUserParty.Put("/reservation_fee", UpdateReservationFee)
	reservationUserParty.Delete("/reservation_fee/{ID}", DeleteReservationFee)

	reservationUserParty.Post("/reservation_report", CreateReservationReport)
	reservationUserParty.Put("/reservation_report", UpdateReservationReport)
	reservationUserParty.Delete("/reservation_report/{ID}", DeleteReservationReport)
	reservationUserParty.Get("/reservation_reports", GetReservationReports)

	userParty := api.Party("", RightsMiddleWare(&userHandler))
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

	userParty.Get("/commitments/forecasts", GetCmtForecasts)

	userParty.Get("/beneficiaries", GetBeneficiaries)
	userParty.Get("/beneficiaries/paginated", GetPaginatedBeneficiaries)
	userParty.Get("/beneficiary/{ID}/datas", GetPaginatedBeneficiaryDatas)
	userParty.Get("/beneficiary/{ID}/export", GetExportBeneficiaryDatas)
	userParty.Get("/beneficiary/{ID}/payments", GetBeneficiaryPayments)
	userParty.Get("/beneficiary/{ID}/placements", GetBeneficiaryPlacements)

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
	userParty.Get("/ratios/years", GetPmtRatiosYears)

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

	userParty.Get("/payment_credit_journal", GetAllPaymentCreditJournals)

	userParty.Get("/payment_credits_and_journal", GetPaymentCreditsAndJournal)

	userParty.Get("/placements", GetPlacements)

	userParty.Get("/beneficiary_groups", GetBeneficiaryGroups)
	userParty.Get("/beneficiary_group/{ID}", GetBeneficiaryGroupItems)
	userParty.Get("/beneficiary_group/{ID}/datas", GetPaginatedBeneficiaryGroupDatas)
	userParty.Get("/beneficiary_group/{ID}/export", GetExportBeneficiaryGroupDatas)
	userParty.Get("/beneficiary_group/{ID}/placements", GetBeneficiaryGroupPlacements)

	userParty.Get("/housing_typologies", GetHousingTypologies)

	userParty.Get("/housing_conventions", GetHousingConventions)

	userParty.Get("/housing_comments", GetHousingComments)

	userParty.Get("/housing_transfers", GetHousingTransfers)

	userParty.Get("/convention_types", GetConventionTypes)

	userParty.Get("/reservation_fees/settings", GetReservationFeeSettings)

	userParty.Get("/housing_types", GetHousingTypes)

	userParty.Get("/dif_action_pmt_prev", GetDifActionPaymentPrevisions)
}
