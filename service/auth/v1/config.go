package v1

import (
	"github.com/go-chi/jwtauth/v5"
	"os"
)

var secretKey *string
var TokenAuth *jwtauth.JWTAuth

const defaultSecretKey = "veGmwk4gcgKKRgQwMnw4t*_s"

const secretKeyEnvKey = "SECRET_KEY"

func init() {
	// Config from env is read by a separate function so that we could call it once
	// again in tests, after setting env variable inside the test
	readConfig()
}

// readConfig initializes config if it is not yet set,
// then reads env variable from databaseConnectionEnvKey
// if the variable is set and config has just been initialized
func readConfig() {
	// If databaseConnection has not been set yet, initialize it
	if secretKey == nil {
		defaultKey := defaultSecretKey
		secretKey = &defaultKey
	}

	// Here we check for two things:
	// 1. There is secretKeyEnvKey environment variable
	// 2. secretKey equals defaultSecretKey, which means
	// that it has just been initialized and was not overwritten previously
	// by flag.Parse(). This is required for tests.
	if key, ok := os.LookupEnv(secretKeyEnvKey); ok && *secretKey == defaultSecretKey {
		secretKey = &key
	}

	TokenAuth = jwtauth.New("HS256", []byte(*secretKey), nil)
}
