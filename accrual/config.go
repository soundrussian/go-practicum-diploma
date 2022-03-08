package accrual

import (
	"flag"
	"os"
)

var accrualAddress *string

const defaultAccrualAddress = "http://localhost:1337"

const accrualAddressEnvKey = "ACCRUAL_SYSTEM_ADDRESS"
const accrualAddressFlag = "r"

func init() {
	// Config from env is read by a separate function so that we could call it once
	// again in tests, after setting env variable inside the test
	readConfig()

	flag.StringVar(accrualAddress, accrualAddressFlag, *accrualAddress, "accrual system address")
}

// readConfig initializes config if it is not yet set,
// then reads env variable from accrualAddressEnvKey
// if the variable is set and config has just been initialized
func readConfig() {
	// If accrualAddress has not been set yet, initialize it
	if accrualAddress == nil {
		defaultAddress := defaultAccrualAddress
		accrualAddress = &defaultAddress
	}

	// Here we check for two things:
	// 1. There is accrualAddressEnvKey environment variable
	// 2. accrualAddress equals defaultAccrualAddress, which means
	// that it has just been initialized and was not overwritten previously
	// by flag.Parse(). This is required for tests.
	if address, ok := os.LookupEnv(accrualAddressEnvKey); ok && *accrualAddress == defaultAccrualAddress {
		accrualAddress = &address
	}
}
