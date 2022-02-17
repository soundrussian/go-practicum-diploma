package api

import (
	"flag"
	"os"
)

// Config contains API settings needed to run gophermart server
type Config struct {
	// RunAddress contains host and port to run server on
	RunAddress string
}

// config is a pointer to API config.
// It is a pointer to simplify testing setting config from env and command-line flags
var config *Config

const defaultRunAddress = "localhost:8080"
const runAddressEnvKey = "RUN_ADDRESS"
const runAddressFlagName = "a"

func init() {
	// Config from env is read by a separate function so that we could call it once
	// again in tests, after setting env variable inside the test
	readConfig()

	flag.StringVar(&config.RunAddress, runAddressFlagName, config.RunAddress, "address to run gophermart on")
}

// readConfig initializes config if it is not yet set,
// then reads env variable from runAddressEnvKey
// if the variable is set and config has just been initialized
func readConfig() {
	// If config has not been set yet, initialize it
	if config == nil {
		config = &Config{RunAddress: defaultRunAddress}
	}

	// Here we check for two things:
	// 1. There is runAddressEnvKey environment variable
	// 2. config.RunAddress equals defaultRunAddress, which means
	// that it has just been initialized and was not overwritten previously
	// by flag.Parse(). This is required for tests.
	if address, ok := os.LookupEnv(runAddressEnvKey); ok && config.RunAddress == defaultRunAddress {
		config.RunAddress = address
	}
}
