package handlerMessage

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	// ErrMandatoryParameterMissing  error in case of mandatory parameter is missing in the request
	ErrInputParameters = errors.New("Error on input request")
	// ErrResourceAlreadyExist  error in case the resource already exist
	ErrResourceAlreadyExist = errors.New("Resource Already Exist")
	// ErrInternalError  error in case internal server erro
	ErrInternalError = errors.New("Internal Error")
	// ErrSecurityToken  error in case the provided mt is not valid
	ErrSecurityToken = errors.New("Mt token not valid")
	// ErrServicePokemon  error in case the provided pokemon name is not valid or the service is unavailable
	ErrServicePokemon = errors.New("Error while fetch pokemon info. Please check your connection or insert valid pokemon name")
	// ErrPokemonName  error in case the provided pokemon name is not valid
	ErrPokemonCategory = errors.New("This category is not covered by insurance")
	// ErrQuoteNotFound  error in case doesn't exist an insurance for the pokemon
	ErrQuoteNotFound = errors.New("Doesn't exist a quote for this pokemon")
	// ErrPaymentSystem  error in case doesn't exist an insurance for the pokemon
	ErrPaymentSystem = errors.New("Payment failed")
	// NoErrorResourceCreated in case a new resource is created
	NoErrorResourceCreated = errors.New("Resource Created")
)

// translate the internal error code to http StatusCode and related message
func ToStatusCodeMessage(err error) (int, string) {
	switch err {
	case ErrInputParameters:
		return http.StatusBadRequest, err.Error()
	case ErrResourceAlreadyExist:
		return http.StatusConflict, err.Error()
	case ErrSecurityToken:
		return http.StatusUnauthorized, err.Error()
	case NoErrorResourceCreated:
		return http.StatusCreated, err.Error()
	case ErrServicePokemon, ErrPokemonCategory, ErrPaymentSystem:
		return http.StatusBadRequest, err.Error()
	case ErrQuoteNotFound:
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "Internal Error"
	}
}
