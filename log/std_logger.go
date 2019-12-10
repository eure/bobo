package log

import (
	"log"
)

// StdLogger use standard log package.
type StdLogger struct {
	IsDebug bool
}

// Debugf logging debug information.
func (l *StdLogger) Debugf(prefix, format string, v ...interface{}) {
	if l.IsDebug {
		log.Printf("[DEBUG] ["+prefix+"] "+format, v...)
	}
}

// Infof logging information.
func (*StdLogger) Infof(prefix, format string, v ...interface{}) {
	log.Printf("[INFO] ["+prefix+"] "+format, v...)
}

// Errorf logging error information.
func (*StdLogger) Errorf(prefix, format string, v ...interface{}) {
	log.Printf("[ERROR] ["+prefix+"] "+format, v...)
}
