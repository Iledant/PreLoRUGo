package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kataras/iris"

	// Imported in config to avoid double import
	_ "github.com/lib/pq"

	yaml "gopkg.in/yaml.v2"
)

// AppStage defines the if the application is used for test, development or
// production.
const (
	ProductionStage  = 1
	DevelopmentStage = 2
	TestStage        = 3
)

// PreLoRuGoConf embeddes the configuration options of the application.
// It's structure is designed to match to the yaml config file used by tests
// and development stages.
// When deployed on production, PreLoRuGoConf uses environnement variables.
// Databases field embeddes configurations fields for tests, development and
// production.
// Users are only for tests purpose to check routes protections.
// App is used to define stage (production or not), logger configuration and
// token file name for persisting tokens on server reload.
type PreLoRuGoConf struct {
	Databases Databases
	Users     Users
	App       App
}

// Credentials embeddes email and password for a user.
type Credentials struct {
	Email, Password, Token string
}

// Users includes users credentials for test purposes.
type Users struct {
	SuperAdmin       Credentials `yaml:"superadmin"`
	Admin            Credentials
	User             Credentials
	CoproUser        Credentials `yaml:"coprouser"`
	RenewProjectUser Credentials `yaml:"renewprojectuser"`
	HousingUser      Credentials `yaml:"housinguser"`
}

// DBConf includes all informations for connecting to a database.
type DBConf struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

// Databases includes databases settings for production, development and tests.
type Databases struct {
	Prod        DBConf
	Development DBConf
	Test        DBConf
}

// App defines global configuration fields for the application (stage, log and
// token file name).
type App struct {
	Stage         int    `yaml:"stage"`
	LogFileName   string `yaml:"logfilename"`
	LoggerLevel   string `yaml:"loggerlevel"`
	TokenFileName string `yaml:"tokenfilename"`
}

var config *PreLoRuGoConf

func logFileOpen(name string, app *iris.Application) (*os.File, error) {
	logFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	app.Logger().SetOutput(logFile)
	app.Logger().Infof("Fichier log configuré")
	return logFile, err
}

// Get fetches all parameters according to tne context : if proper environment
// variables are set, assumes beeing in prod, otherwise read the config.yml file
func (p *PreLoRuGoConf) Get(app *iris.Application) (logFile *os.File, err error) {
	if config != nil {
		p = config
		return nil, nil
	}
	// Configure the log file as first step to catch all messages
	p.App.LogFileName = os.Getenv("LOG_FILE_NAME")
	if p.App.LogFileName != "" {
		logFile, err = logFileOpen(p.App.LogFileName, app)
		if err != nil {
			return nil, err
		}
	}

	// Check if RDS environment variables are set
	name := os.Getenv("RDS_DB_NAME")
	host := os.Getenv("RDS_HOSTNAME")
	port := os.Getenv("RDS_PORT")
	username := os.Getenv("RDS_USERNAME")
	password := os.Getenv("RDS_PASSWORD")

	if name != "" && host != "" && port != "" && username != "" && password != "" {
		app.Logger().Infof("Utilisation des variables d'environnement")
		p.Databases.Prod.Name = name
		p.Databases.Prod.Host = host
		p.Databases.Prod.Port = port
		p.Databases.Prod.UserName = username
		p.Databases.Prod.Password = password
		p.App.TokenFileName = os.Getenv("TOKEN_FILE_NAME")
		p.Users.SuperAdmin.Email = os.Getenv("SUPERADMIN_EMAIL")
		p.Users.SuperAdmin.Password = os.Getenv("SUPERADMIN_PWD")
		p.App.Stage = ProductionStage
		app.Logger().SetLevel("info")
		return logFile, nil
	}

	// Otherwise try to fetch configuration through yml configuration file
	cfgFile, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		// Check an other location for test purposes
		cfgFile, err = ioutil.ReadFile("config.yml")
		if err != nil {
			return logFile, fmt.Errorf("Erreur de lecture de config.yml : %v", err)
		}
	}
	if err = yaml.Unmarshal(cfgFile, p); err != nil {
		return logFile, fmt.Errorf("Erreur de  décodage de config.yml : %v", err)
	}
	if p.App.LoggerLevel != "" {
		app.Logger().SetLevel(p.App.LoggerLevel)
	}
	if logFile == nil && p.App.LogFileName != "" {
		logFile, err = logFileOpen(p.App.LogFileName, app)
	}
	app.Logger().Infof("Utilisation de config.yml")
	return logFile, err
}
