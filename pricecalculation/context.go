package pricecalculation

import "github.com/rs/zerolog"

type Context struct {
	logger                    zerolog.Logger
	resetVisitHistories       ForResettingVisitHistories
	InitCounter               uint
}

func NewContext(logger zerolog.Logger, resetVisitHistories ForResettingVisitHistories) *Context {
	return &Context{logger: logger, resetVisitHistories: resetVisitHistories}
}

func (context *Context) Initialize() {
	context.resetVisitHistories()
	context.InitCounter++

	context.logger.Debug().Msgf("context init count: %d", context.InitCounter)
}
