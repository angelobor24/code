package handlerMessage

import (
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func TestNewStorageImpl(t *testing.T) {
	code, err := ToStatusCodeMessage(ErrInputParameters)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, err, ErrInputParameters.Error())

	code, err = ToStatusCodeMessage(ErrResourceAlreadyExist)
	assert.Equal(t, code, http.StatusConflict)
	assert.Equal(t, err, ErrResourceAlreadyExist.Error())

	code, err = ToStatusCodeMessage(ErrSecurityToken)
	assert.Equal(t, code, http.StatusUnauthorized)
	assert.Equal(t, err, ErrSecurityToken.Error())

	code, err = ToStatusCodeMessage(NoErrorResourceCreated)
	assert.Equal(t, code, http.StatusCreated)
	assert.Equal(t, err, NoErrorResourceCreated.Error())

	code, err = ToStatusCodeMessage(ErrServicePokemon)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, err, ErrServicePokemon.Error())

	code, err = ToStatusCodeMessage(ErrQuoteNotFound)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, err, ErrQuoteNotFound.Error())

	code, err = ToStatusCodeMessage(ErrPaymentSystem)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, err, ErrPaymentSystem.Error())

	code, err = ToStatusCodeMessage(http.ErrAbortHandler)
	assert.Equal(t, code, http.StatusInternalServerError)
	assert.Equal(t, err, "Internal Error")
}
