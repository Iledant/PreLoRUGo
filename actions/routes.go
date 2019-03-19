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
	adminParty.Post("/housing", CreateHousing)

	adminParty.Post("/housing", CreateHousing)
	adminParty.Put("/housing", UpdateHousing)
	adminParty.Delete("/housing/{ID}", DeleteHousing)
	adminParty.Post("/housings", BatchHousings)

	adminParty.Post("/commitments", BatchCommitments)
	adminParty.Post("/commitments/link", SetCommitmentLink)

	adminParty.Post("/payments", BatchPayments)

	adminParty.Post("/budget_sector", CreateBudgetSector)
	adminParty.Put("/budget_sector", UpdateBudgetSector)
	adminParty.Delete("/budget_sector/{ID}", DeleteBudgetSector)

	userParty := api.Party("", ActiveMiddleware)
	userParty.Post("/user/password", ChangeUserPwd)
	userParty.Get("/user/logout", Logout)
	userParty.Get("/budget_actions", GetBudgetActions)
	userParty.Get("/copro", GetCopros)
	userParty.Get("/renew_projects", GetRenewProjects)

	userParty.Get("/housings", GetHousings)

	userParty.Get("/commitments", GetCommitments)

	userParty.Get("/beneficiaries", GetBeneficiaries)

	userParty.Get("/payments", GetPayments)

	userParty.Get("/budget_sector/{ID}", GetBudgetSector)
	userParty.Get("/budget_sectors", GetBudgetSectors)
}

// setDBMiddleware return a middleware to add db to context values
func setDBMiddleware(db *sql.DB) func(iris.Context) {
	return func(ctx iris.Context) {
		ctx.Values().Set("db", db)
		ctx.Next()
	}
}
