package enum

type InvoicePaymentType string

const (
	InvoicePaymentType_CASH   = InvoicePaymentType("CASH")
	InvoicePaymentType_CREDIT = InvoicePaymentType("CREDIT")
)

var validInvoicePaymentType = []InvoicePaymentType{
	InvoicePaymentType_CASH,
	InvoicePaymentType_CREDIT,
}

func (InvoicePaymentType InvoicePaymentType) String() string {
	return string(InvoicePaymentType)
}

func (InvoicePaymentType InvoicePaymentType) IsValid() bool {
	for _, item := range validInvoicePaymentType {
		if item == InvoicePaymentType {
			return true
		}
	}
	return false
}
