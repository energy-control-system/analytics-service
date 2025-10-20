package cron

import "github.com/sunshineOfficial/golib/golog"

type logger struct {
	log golog.Logger
}

func (l logger) Debug(msg string, args ...any) {
	l.log.Debugf(msg, args...)
}

func (l logger) Error(msg string, args ...any) {
	l.log.Errorf(msg, args...)
}

func (l logger) Info(msg string, args ...any) {
	l.log.Debugf(msg, args...)
}

func (l logger) Warn(msg string, args ...any) {
	l.log.Errorf(msg, args...)
}
