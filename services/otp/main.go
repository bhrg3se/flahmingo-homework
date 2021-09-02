package main

import (
	"flag"
	"github.com/bhrg3se/flahmingo-homework/utils"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	path := flag.String("c", "/etc/flahmingo", "config file location")
	writeToFile := flag.Bool("f", false, "write logs to file")
	flag.Parse()

	config := utils.ParseConfig(*path)
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
	//logrus.SetReportCaller(true)

	if *writeToFile {
		f, err := os.OpenFile("/var/log/flahmingo/otp.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	}

	startService(config)

}
