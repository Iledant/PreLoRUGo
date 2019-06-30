package main

import (
	stdContext "context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Iledant/PreLoRUGo/actions"
	"github.com/Iledant/PreLoRUGo/config"
	"github.com/kataras/iris"
)

func run(app *iris.Application, cfg *config.PreLoRuGoConf) error {
	if cfg.App.Stage == config.ProductionStage {
		domain := os.Getenv("APP_DOMAIN")
		if domain == "" {
			app.Logger().Error("Variable d'environnement APP_DOMAIN vide")
			return fmt.Errorf("Mauvaise configuration des variables d'environnement")
		}
		addr := os.Getenv("APP_ADDR")
		if addr == "" {
			app.Logger().Error("Variable d'environnement APP_ADDR vide")
			return fmt.Errorf("Mauvaise configuration des variables d'environnement")
		}
		email := os.Getenv("DOMAIN_OWNER_EMAIL")
		if email == "" {
			app.Logger().Error("Variable d'environnement DOMAIN_OWNER_EMAIL vide")
			return fmt.Errorf("Mauvaise configuration des variables d'environnement")
		}
		crtDir := os.Getenv("CRT_DIR")
		if crtDir == "" {
			app.Logger().Error("Variable d'environnement CRT_DIR vide")
			return fmt.Errorf("Mauvaise configuration des variables d'environnement")
		}
		return app.NewHost(&http.Server{Addr: ":443"}).
			ListenAndServeAutoTLS(domain, email, crtDir)
	}
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
	err = run(app, &cfg)
	app.Logger().Fatalf("Erreur de serveur run %v", err)
}
