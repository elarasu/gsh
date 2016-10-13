package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

// setup log levels
func init() {
	// set this if environment is production
	//	log.SetFormatter(&log.JSONFormatter{})

	// set this for non prod
	log.SetFormatter(&log.TextFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("debug statement")
	log.Info("info text")
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
}
