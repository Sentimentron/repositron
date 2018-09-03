package database

import "os"

func createDatabaseForTesting(tempPath string) error {
	os.Remove(tempPath)
	return CreateDatabaseIfNotExists(tempPath)
}
