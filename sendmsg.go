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
	switch os.Args[1] {
	case "simple":
		simpleCommand.Parse(os.Args[2:])
		cfglocation = simpleFrontendCmd.SimpleCFGLocation
		break
	case "icingaHost":
		icingaHostCommand.Parse(os.Args[2:])
		cfglocation = icingaHostCmd.SimpleCFGLocation
		break
	case "icingaService":
		icingaServiceCommand.Parse(os.Args[2:])
		cfglocation = icingaServiceCmd.SimpleCFGLocation
		break
	default:
		flag.PrintDefaults()
		return
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
			"Body": msg.Body,
			"Title": msg.Body_title,
		}).Fatal("Either title or body are missing!")
	}

	switch cfg.Backend {
	case "slack":
		send_with_slack(msg, cfg)
	}
}
