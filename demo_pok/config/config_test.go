package config

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestGetOrDefault(t *testing.T) {
	envVarName := "DUMMY_ENV"
	defaultEnvValue := "testDefaultDummyEnv"
	assert.Equal(t, getOrDefault(envVarName, defaultEnvValue), defaultEnvValue)
	newEnvValue := "new-value-for-dummy-env"
	os.Setenv(envVarName, newEnvValue)
	assert.Equal(t, getOrDefault(envVarName, defaultEnvValue), newEnvValue)
	os.Unsetenv(envVarName)
}

func TestFlagFmt(t *testing.T) {
	str := "INT_DUMMY_ENV"
	expected := "int-dummy-env"
	actual := flagFmt(str)
	assert.Equal(t, actual, expected)
}

func TestEnvFmt(t *testing.T) {
	str := "test-FOO-Env"
	expected := "TEST_FOO_ENV"
	actual := envFmt(str)
	assert.Equal(t, actual, expected)
}

func TestGetDefaultPaymentBehaviour(t *testing.T) {
	paymentBehaviour := GetDefaultPaymentBehaviour()
	assert.Equal(t, paymentBehaviour, defaultPaymentBehaviour)
}

func TestGetPayBehaviour(t *testing.T) {
	paymentBehaviour := GetPayBehaviour()
	assert.Equal(t, paymentBehaviour, defaultPaymentBehaviour)
	os.Setenv("PAYMENT_BEHAVIOUR", "rejected")
	paymentBehaviour = GetPayBehaviour()
	assert.Equal(t, paymentBehaviour, "rejected")
	os.Unsetenv("PAYMENT_BEHAVIOUR")
}
