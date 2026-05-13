package pricecalculation

import (
	"implementing-ddd-in-go/pricecalculation/domain"

	"github.com/rs/zerolog"
)

type Context struct {
	logger         zerolog.Logger
	GetVisitorByID ForGettingVisitorByID
	VisitHistory   *domain.VisitHistory
	InitCounter    uint
}

func NewContext(logger zerolog.Logger) *Context {
	return &Context{logger: logger, VisitHistory: domain.NewVisitHistory()}
}

func (context *Context) Initialize(getVisitorByID ForGettingVisitorByID) {
	context.GetVisitorByID = getVisitorByID
	context.VisitHistory.Reset()
	context.InitCounter++

	context.logger.Debug().Msgf("context init count: %d", context.InitCounter)
}
