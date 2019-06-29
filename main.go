package main

import (
	stdContext "context"
	"time"

	"github.com/Iledant/PreLoRUGo/actions"
	"github.com/Iledant/PreLoRUGo/config"
	"github.com/kataras/iris"
)

func run(app *iris.Application, cfg *config.PreLoRuGoConf) error {
	// if cfg.App.Prod {
	// 	domain := os.Getenv("APP_DOMAIN")
	// 	addr := os.Getenv("APP_ADDR")
	// 	email := os.Getenv("DOMAIN_OWNER_EMAIL")
	// 	crtDir := os.Getenv("CRT_DIR")
	// 	return app.NewHost(&http.Server{Addr: addr}).
	// 		ListenAndServeAutoTLS(domain, email, crtDir)
	// }
	return app.Run(iris.Addr(":5000"))
}

func main() {
	app := iris.New().Configure(
		iris.WithConfiguration(iris.Configuration{DisablePathCorrection: true}))

	var cfg config.PreLoRuGoConf
	logFile, err := cfg.Get(app)
	if logFile != nil {
		defer logFile.Close()
	}
	if err != nil {
		app.Logger().Fatalf("Configuration : %v", err)
	}

	db, err := config.InitDatabase(&cfg, false, true)
	if err != nil {
		app.Logger().Fatalf("Initialisation de la base de données : %v", err)
	}
	app.Logger().Infof("Base de données connectée et initialisée")
	defer db.Close()

	actions.SetRoutes(app, db)
	app.StaticWeb("/", "./dist")
	app.Logger().Infof("Routes et serveur statique configurés")

	// Configure tokens recover and autosave on stop
	if cfg.App.TokenFileName != "" {
		actions.TokenRecover(cfg.App.TokenFileName)
		iris.RegisterOnInterrupt(func() {
			timeout := 2 * time.Second
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
			defer cancel()
			actions.TokenSave(cfg.App.TokenFileName)
			app.Shutdown(ctx)
		})
		app.Logger().Infof("Fichier de sauvegarde des tokens configuré")
	}
	// Run application according to application stage
	run(app, &cfg)
}
