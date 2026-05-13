package pricecalculation

import (
	"implementing-ddd-in-go/pricecalculation/domain"

	"github.com/rs/zerolog"
)

type forGettingVisitorByID func(id string) (domain.ExternalVisitor, error)

type Context struct {
	logger zerolog.Logger
	GetVisitorByID forGettingVisitorByID
	InitCounter    uint
}

func NewContext(logger zerolog.Logger) *Context {
	return &Context{logger: logger}
}

func (context *Context) Initialize(getVisitorByID forGettingVisitorByID) {
	context.GetVisitorByID = getVisitorByID
	context.InitCounter++

	context.logger.Debug().Msgf("context init count: %d", context.InitCounter)
}

