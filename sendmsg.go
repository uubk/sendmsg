package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var gitversion = "undefined"

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
	// 'Simple' frontend - that is specifying each field manually
	simpleCommand := flag.NewFlagSet("simple", flag.ExitOnError)
	simpleFrontendCmd := simpleCmdData{}
	simpleFrontendCmd.Init(simpleCommand)

	// 'Icinga2' frontend - use this as a drop-in replacement for mails
	icingaHostCommand := flag.NewFlagSet("icingaHost", flag.ExitOnError)
	icingaServiceCommand := flag.NewFlagSet("icingaService", flag.ExitOnError)
	icingaHostCmd := icingaCmd{}
	icingaHostCmd.Init(icingaHostCommand, false)
	icingaServiceCmd := icingaCmd{}
	icingaServiceCmd.Init(icingaServiceCommand, true)

	if len(os.Args) < 2 {
		logrus.Fatal("Please specify a command [simple, icingaHost, icingaService]")
		flag.PrintDefaults()
	}

	var cfglocation *string
	var debug *bool
	switch os.Args[1] {
	case "simple":
		simpleCommand.Parse(os.Args[2:])
		cfglocation = simpleFrontendCmd.CFGLocation
		debug = simpleFrontendCmd.Debug
		break
	case "icingaHost":
		icingaHostCommand.Parse(os.Args[2:])
		cfglocation = icingaHostCmd.CFGLocation
		debug = icingaHostCmd.Debug
		break
	case "icingaService":
		icingaServiceCommand.Parse(os.Args[2:])
		cfglocation = icingaServiceCmd.CFGLocation
		debug = icingaServiceCmd.Debug
		break
	default:
		flag.PrintDefaults()
		return
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg := Config{}
	if contents, err := ioutil.ReadFile(*cfglocation); err == nil {
		err = yaml.UnmarshalStrict(contents, &cfg)
		if err != nil {
			logrus.WithError(err).Fatal("Couldn't parse the configuration file!")
		}
	} else {
		logrus.WithError(err).Fatal("Couldn't read the configuration file!")
	}

	logrus.WithFields(logrus.Fields{
		"version":  gitversion,
		"frontend": os.Args[1],
		"backend":  cfg.Backend,
	}).Info("Starting sendmsg")

	var msg Message
	if simpleCommand.Parsed() {
		msg = simpleFrontendCmd.Parse()
	}
	if icingaHostCommand.Parsed() {
		msg = icingaHostCmd.Parse()
	}
	if icingaServiceCommand.Parsed() {
		msg = icingaServiceCmd.Parse()
	}

	if msg.Body == "" || msg.Body_title == "" {
		logrus.WithFields(logrus.Fields{
			"Body":  msg.Body,
			"Title": msg.Body_title,
		}).Fatal("Either title or body are missing!")
	}

	switch cfg.Backend {
	case "slack":
		send_with_slack(msg, cfg)
	default:
		logrus.WithField("backend", cfg.Backend).Fatal("Unkown backend!")
	}
}
