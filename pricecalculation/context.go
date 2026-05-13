package pricecalculation

import (
	"github.com/rs/zerolog"
)

type Context struct {
	logger         zerolog.Logger
	GetVisitorByID ForGettingVisitorByID
	InitCounter    uint
}

func NewContext(logger zerolog.Logger) *Context {
	return &Context{logger: logger}
}

func (context *Context) Initialize(getVisitorByID ForGettingVisitorByID) {
	context.GetVisitorByID = getVisitorByID
	context.InitCounter++

	context.logger.Debug().Msgf("context init count: %d", context.InitCounter)
}
