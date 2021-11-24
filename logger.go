package paasport

import "log"

// Logger logger define
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct{}

var defaultLogger logger

func (l logger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func (l logger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l logger) Warnf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func (l logger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
