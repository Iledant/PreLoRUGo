package config

import (
	"errors"
	"io/ioutil"
	"os"

	// Imported in config to avoid double import
	_ "github.com/lib/pq"

	yaml "gopkg.in/yaml.v2"
)

// PreLoRuGoConf embeddes the configuration to decde yaml file
type PreLoRuGoConf struct {
	Databases Databases
	Users     Users
	App       App
}

// Users includes users credentials for test purposes.
type Users struct {
	Admin            Credentials
	User             Credentials
	CoproUser        Credentials `yaml:"coprouser"`
	RenewProjectUser Credentials `yaml:"renewprojectuser"`
	HousingUser      Credentials `yaml:"housinguser"`
}

// Databases includes the 3 databases settings for production, development and tests.
type Databases struct {
	Prod        DBConf
	Development DBConf
	Test        DBConf
}

// App defines global values for the application
type App struct {
	Prod          bool   `yaml:"prod"`
	LogFileName   string `yaml:"logfilename"`
	LoggerLevel   string `yaml:"loggerlevel"`
	TokenFileName string `yaml:"tokenfilename"`
}

// DBConf includes all informations for connecting to a database.
type DBConf struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

// Credentials keep email ans password for a user.
type Credentials struct {
	Email, Password, Token string
}

var config *PreLoRuGoConf

// Get fetches all parameters according to tne context : if proper environment
// variables are set, assumes beeing in prod, otherwise read the config.yml file
func (p *PreLoRuGoConf) Get() error {
	if config == nil {
		// Check if RDS environment variables are set
		name := os.Getenv("RDS_DB_NAME")
		host := os.Getenv("RDS_HOSTNAME")
		port := os.Getenv("RDS_PORT")
		username := os.Getenv("RDS_USERNAME")
		password := os.Getenv("RDS_PASSWORD")

		if name != "" && host != "" && port != "" && username != "" && password != "" {
			p = &PreLoRuGoConf{Databases: Databases{Prod: DBConf{
				Name:     name,
				Host:     host,
				Port:     port,
				UserName: username,
				Password: password}}}
			p.App.TokenFileName = os.Getenv("TOKEN_FILE_NAME")
			p.App.LogFileName = os.Getenv("LOG_FILE_NAME")
			p.App.Prod = true
			p.App.LoggerLevel = "warn"
			return nil
		}
		// Otherwise use database.yml
		cfgFile, err := ioutil.ReadFile("../config.yml")
		if err != nil {
			// Try to read directly
			cfgFile, err = ioutil.ReadFile("config.yml")
			if err != nil {
				return errors.New("Erreur lors de la lecture de config.yml : " + err.Error())
			}
		}
		if err = yaml.Unmarshal(cfgFile, p); err != nil {
			return errors.New("Erreur lors du décodage de config.yml : " + err.Error())
		}
	} else {
		p = config
	}
	return nil
}
