package domain

type InvoicingPolicy func(PriceCalculated) bool

var BusinessCustomersRequireInvoice InvoicingPolicy = func(e PriceCalculated) bool {
	return e.CustomerType == customerTypeBusiness
}
