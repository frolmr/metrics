package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	SugaredLogger zap.SugaredLogger
}

func NewLogger() (*Logger, error) {
	l, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	return &Logger{
		SugaredLogger: *l.Sugar(),
	}, nil
}
