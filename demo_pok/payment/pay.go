package payment

import "demo/config"

type Payment interface {
	Pay() bool
}

type PaymentImpl struct {
}

type MockPaymentImpl struct {
	expectedEsito bool
}

func NewPaymentImpl() Payment {
	paymentImpl := PaymentImpl{}
	return &paymentImpl
}

// return the payment result operation
func (paymentImpl *PaymentImpl) Pay() bool {
	return config.GetPayBehaviour() == config.GetDefaultPaymentBehaviour()
}

// return a mock of the payment interface
func MockNewPaymentImpl(esito bool) Payment {
	paymentImpl := MockPaymentImpl{expectedEsito: esito}
	return &paymentImpl
}

// mock implementation of Pay() interface function
func (paymentImpl *MockPaymentImpl) Pay() bool {
	return paymentImpl.expectedEsito
}
