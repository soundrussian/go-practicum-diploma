package v1

import (
	"flag"
	"os"
)

var databaseConnection *string

const defaultDatabaseConnection = "port=5432 host=localhost user=postgres password=postgres dbname=my_database sslmode=disable"

const databaseConnectionEnvKey = "DATABASE_URI"
const databaseConnectionFlag = "d"

func init() {
	// Config from env is read by a separate function so that we could call it once
	// again in tests, after setting env variable inside the test
	readConfig()

	flag.StringVar(databaseConnection, databaseConnectionFlag, *databaseConnection, "database connection")
}

// readConfig initializes config if it is not yet set,
// then reads env variable from databaseConnectionEnvKey
// if the variable is set and config has just been initialized
func readConfig() {
	// If databaseConnection has not been set yet, initialize it
	if databaseConnection == nil {
		defaultDSN := defaultDatabaseConnection
		databaseConnection = &defaultDSN
	}

	// Here we check for two things:
	// 1. There is databaseConnectionEnvKey environment variable
	// 2. databaseConnection equals defaultRunAddress, which means
	// that it has just been initialized and was not overwritten previously
	// by flag.Parse(). This is required for tests.
	if dsn, ok := os.LookupEnv(databaseConnectionEnvKey); ok && *databaseConnection == defaultDatabaseConnection {
		databaseConnection = &dsn
	}
}
