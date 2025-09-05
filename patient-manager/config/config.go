package config

const (
	LOG_FILE                    = "ePrometna.log"
	LOG_FILE_MAX_SIZE           = 2
	LOG_FILE_MAX_AGE            = 30
	LOG_FILE_MAX_BACKUPS        = 0
	LOG_FOLDER           string = "./log"
	TMP_FOLDER           string = "./tmp"
)

// AppConfig is struct that contains basic app configuration variables
var AppConfig *AppConfiguration = nil

type AppConfiguration struct {
	Env                environment
	Port               int
	DbConnection       string
	AccessKey          string
	RefreshKey         string
	MIOEndpoint        string
	MIOAccessKeyID     string
	MIOSecretAccessKey string
	UseSSL             bool
}

type environment = string

const (
	Dev  = "dev"
	Test = "test"
	Prod = "prod"
)
