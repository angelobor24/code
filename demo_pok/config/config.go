package config

import (
	"flag"
	"os"
	"strings"
)

const (
	paymentBehaviour        = "PAYMENT_BEHAVIOUR"
	defaultPaymentBehaviour = "accept"
)

// retrieve the expected behaviour to simulate the payment
func GetPayBehaviour() string {
	return getOrDefault(paymentBehaviour, defaultPaymentBehaviour)
}

// extract the ENV variable
func getOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(envFmt(key))
	if exists {
		return value
	}
	f := flag.Lookup(flagFmt(key))
	if f == nil {
		return defaultValue
	}
	return f.Value.String()
}

// read ENV from environment
func envFmt(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, "-", "_"))
}

// read ENV from command line
func flagFmt(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, "_", "-"))
}

// return defaultt behaviour
func GetDefaultPaymentBehaviour() string {
	return defaultPaymentBehaviour
}
