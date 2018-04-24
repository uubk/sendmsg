package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
)

func main() {
	// 'Simple' frontend - that is specifying each field manually
	simpleCommand := flag.NewFlagSet("simple", flag.ExitOnError)
	simpleFrontendCmd := simpleCmdData{}
	simpleFrontendCmd.Init(simpleCommand)

	// 'Icinga2' frontend - use this as a drop-in replacement for mails
	icingaHostCommand := flag.NewFlagSet("icingaHost", flag.ExitOnError)
	icingaHostCmd := icingaHostCmd{}
	icingaHostCmd.Init(icingaHostCommand)

	if len(os.Args) < 2 {
		fmt.Print("Please use one of [simple, icingaHost]\n")
		return
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
	default:
		flag.PrintDefaults()
		return
	}

	cfg := Config{}
	if contents, err := ioutil.ReadFile(*cfglocation); err == nil {
		err = yaml.UnmarshalStrict(contents, &cfg)
		if err != nil {
			fmt.Println("ERROR: Couldn't parse the configuration file!", err)
			os.Exit(-1)
		}
	} else {
		fmt.Println("ERROR: Couldn't read the configuration file!", err)
		os.Exit(-1)
	}

	var msg Message
	if simpleCommand.Parsed() {
		msg = simpleFrontendCmd.Parse()
	}
	if icingaHostCommand.Parsed() {
		msg = icingaHostCmd.Parse()
	}

	if msg.Body == "" || msg.Body_title == "" {
		fmt.Println("ERROR: Either title or body are missing. Aborting now!")
		return
	}

	switch cfg.Backend {
	case "slack":
		send_with_slack(msg, cfg, "simple", "")
	}
}
