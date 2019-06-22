package main

import (
	stdContext "context"
	"log"
	"time"

	"github.com/Iledant/PreLoRUGo/actions"
	"github.com/Iledant/PreLoRUGo/config"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New().Configure(
		iris.WithConfiguration(iris.Configuration{DisablePathCorrection: true}))

	var cfg config.PreLoRuGoConf
	if err := cfg.Get(); err != nil {
		log.Fatalf("Récupération de la configuration : %v", err)
	}

	db, err := config.InitDatabase(&cfg.Databases.Development, false, true)
	if err != nil {
		log.Fatalf("Initialisation de la base de données : %v", err)
	}
	defer db.Close()
	actions.SetRoutes(app, db)
	if cfg.App.LoggerLevel != "" {
		app.Logger().SetLevel(cfg.App.LoggerLevel)
	}
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
	}
	// Use port 5000 as Elastic beanstalk uses it by default
	app.Run(iris.Addr(":5000"), iris.WithoutInterruptHandler)
}
