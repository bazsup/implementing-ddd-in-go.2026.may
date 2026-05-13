package pricecalculation

import "implementing-ddd-in-go/pricecalculation/domain"

type forGettingVisitorByID func(id string) (domain.ExternalVisitor, error)

type Context struct {
	GetVisitorByID forGettingVisitorByID
}

func NewContext() Context {
	return Context{}
}
