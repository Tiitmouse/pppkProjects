package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Variables

var MongoDb_Conn_String = ""
var Endpoint = ""
var AccessKeyID = ""
var SecretAccessKey = ""
var UseSSL = false

// Variable names
const _MONGODB_CONN_STRING_NAME = "MONGO_CONNECTION_STRING"
const _ENDPOINT_NAME = "MINIO_ENDPOINT"
const _ACCESS_KEY_NAME = "MINIO_ACCESS_KEY_ID"
const _SECRET_ACCESS_KEY_NAME = "MINIO_SECRET_ACCESS_KEY"
const _USE_SSL_NAME = "MINIO_USE_SSL"

func Load() error {
	zap.S().Debugf("Loading env variables")

	if err := godotenv.Load(); err != nil {
		zap.S().Errorf("Env load err = %+v\n", err)
		zap.S().DPanicf("Can't load config using real env")
	}

	MongoDb_Conn_String = loadString(_MONGODB_CONN_STRING_NAME)

	Endpoint = loadString(_ENDPOINT_NAME)
	AccessKeyID = loadString(_ACCESS_KEY_NAME)
	SecretAccessKey = loadString(_SECRET_ACCESS_KEY_NAME)
	UseSSL = loadBool(_USE_SSL_NAME)

	zap.S().Debugf("Finished loading env variables")
	return nil
}

func loadString(name string) string {
	rez := strings.TrimSpace(os.Getenv(name))
	if rez == "" {
		zap.S().Errorf("Env variable %s is empty", name)
	}
	zap.S().Debugf("Loaded %s = %s", name, rez)
	return rez
}

func loadBool(name string) bool {
	rez := os.Getenv(name)
	if rez == "" {
		zap.S().Errorf("Env variable %s is empty", name)
	}
	return rez == "development"
}
