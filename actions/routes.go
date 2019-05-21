package actions

import (
	"database/sql"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

// SetRoutes initialize all routes for the application
func SetRoutes(app *iris.Application, db *sql.DB) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	crsParty := app.Party("/api", crs).AllowMethods(iris.MethodOptions)

	crsParty.Post("/user/sign_up", setDBMiddleware(db), SignUp)
	crsParty.Post("/user/login", setDBMiddleware(db), Login)
	api := crsParty.Party("", setDBMiddleware(db))

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

	adminParty.Post("/city", CreateCity)
	adminParty.Put("/city", UpdateCity)
	adminParty.Delete("/city/{ID}", DeleteCity)
	adminParty.Post("/cities", BatchCities)

	adminParty.Post("/renew_project_forecast", CreateRenewProjectForecast)
	adminParty.Put("/renew_project_forecast", UpdateRenewProjectForecast)
	adminParty.Delete("/renew_project_forecast/{ID}", DeleteRenewProjectForecast)
	adminParty.Post("/renew_project_forecasts", BatchRenewProjectForecasts)

	adminParty.Post("/copro_forecast", CreateCoproForecast)
	adminParty.Put("/copro_forecast", UpdateCoproForecast)
	adminParty.Delete("/copro_forecast/{ID}", DeleteCoproForecast)
	adminParty.Post("/copro_forecasts", BatchCoproForecasts)

	adminParty.Get("/settings", GetSettings)

	adminParty.Post("/ratios", BatchPmtRatios)
	adminParty.Get("/ratios/years", GetPmtRatiosYears)

	userParty := api.Party("", ActiveMiddleware)
	userParty.Post("/user/password", ChangeUserPwd)
	userParty.Post("/user/logout", Logout)
	userParty.Get("/budget_actions", GetBudgetActions)
	userParty.Get("/copro", GetCopros)
	userParty.Get("/renew_projects", GetRenewProjects)

	userParty.Get("/housings", GetHousings)
	userParty.Get("/housings/paginated", GetPaginatedHousings)

	userParty.Get("/commitments", GetCommitments)
	userParty.Get("/commitments/paginated", GetPaginatedCommitments)
	userParty.Get("/commitments/export", ExportCommitments)

	userParty.Get("/beneficiaries", GetBeneficiaries)
	userParty.Get("/beneficiaries/paginated", GetPaginatedBeneficiaries)
	userParty.Get("/beneficiary/{ID}/datas", GetPaginatedBeneficiaryDatas)

	userParty.Get("/payments", GetPayments)
	userParty.Get("/payments/paginated", GetPaginatedPayments)
	userParty.Get("/payments/export", GetExportedPayments)

	userParty.Get("/budget_sectors", GetBudgetSectors)
	userParty.Get("/budget_sector/{ID}", GetBudgetSector)
	userParty.Get("/community/{ID}", GetCommunity)
	userParty.Get("/communities", GetCommunities)

	userParty.Get("/commission/{ID}", GetCommission)
	userParty.Get("/commissions", GetCommissions)

	userParty.Get("/city/{ID}", GetCity)
	userParty.Get("/cities", GetCities)
	userParty.Get("/cities/paginated", GetPaginatedCities)

	userParty.Get("/renew_project_forecast/{ID}", GetRenewProjectForecast)
	userParty.Get("/renew_project_forecasts", GetRenewProjectForecasts)

	userParty.Get("/copro_forecast/{ID}", GetCoproForecast)
	userParty.Get("/copro_forecasts", GetCoproForecasts)

	userParty.Get("/home", GetHome)

	userParty.Get("/ratios", GetPmtRatios)
}

// setDBMiddleware return a middleware to add db to context values
func setDBMiddleware(db *sql.DB) func(iris.Context) {
	return func(ctx iris.Context) {
		ctx.Values().Set("db", db)
		ctx.Next()
	}
}
