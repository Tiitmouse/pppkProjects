// Package app preforms basic app functions like setup, loading config and global definitions
package app

import (
	"data-managment/util/bucket"
	"data-managment/util/env"
	"data-managment/util/repo"
	"fmt"

	"go.uber.org/zap"
)

// Setup will preform app setup or panic of it fails
// Can only be called once
func Setup() {
	// Logger setup
	{
		var err error

		if Build == BuildDev {
			err = devLoggerSetup()
			if err != nil {
				fmt.Printf("err: %v\n", err)
				panic("Failed to setup logger")
			}
		} else {
			err = prodLoggerSetup()
			if err != nil {
				fmt.Printf("err: %v\n", err)
				panic("Failed to setup logger")
			}
		}
	}

	// Load env
	{
		err := env.Load()
		if err != nil {
			zap.S().Panicf("Failed to load env variables")
		}

	}

	// Minio setup
	{
		err := bucket.Setup()
		if err != nil {
			zap.S().Panicf("Failed to setup Minio bucket, err = %+v", err)
		}
	}

	// Minio setup
	{
		err := repo.Setup()
		if err != nil {
			zap.S().Panicf("Failed to setup MongoDb repo, err = %+v", err)
		}
	}
}
