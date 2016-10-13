package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// setup commandline params
var (
	debug   = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout = kingpin.Flag("timeout", "Timeout waiting for ping.").Default("5s").OverrideDefaultFromEnvar("PING_TIMEOUT").Short('t').Duration()
	ip      = kingpin.Arg("ip", "IP address to ping.").Required().IP()
	count   = kingpin.Arg("count", "Number of packets to send").Int()
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

	// command line params
	kingpin.Version("0.0.1")
	kingpin.Parse()
	fmt.Printf("Would ping: %s with timeout %s and count %d", *ip, *timeout, *count)
}
