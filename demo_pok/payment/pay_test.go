package payment

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestPay(t *testing.T) {
	// test mock interface
	paymentInterfaceMock := MockNewPaymentImpl(true)
	assert.Equal(t, paymentInterfaceMock.Pay(), true)
	// test pay interface
	paymentInterface := NewPaymentImpl()
	assert.Equal(t, paymentInterface.Pay(), true)
	os.Setenv("PAYMENT_BEHAVIOUR", "rejected")
	assert.Equal(t, paymentInterface.Pay(), false)
	os.Unsetenv("PAYMENT_BEHAVIOUR")
}
