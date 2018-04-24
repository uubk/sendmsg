package main

import "flag"

type icingaHostCmd struct {
	SimpleCFGLocation *string
	icHostname *string
	icHostV4 *string
	icHostV6 *string
	icHostState *string
	icHostNotfnType *string
	icHostCheckOutput *string
	icHostNotfnUser *string
	icHostNotfnComment *string
	icURL *string
}

func (this *icingaHostCmd) Init(flagSet *flag.FlagSet) {
	this.SimpleCFGLocation = flagSet.String("cfg", "/etc/sendmsg.yml", "Path to sendmsg config")
	this.icHostname = flagSet.String("l", "", "Hostname")
	this.icHostV4 = flagSet.String("4", "", "Host address (v4)")
	this.icHostV6 = flagSet.String("6", "", "Host address (v6)")
	this.icHostState = flagSet.String("s", "", "Host state")
	this.icHostNotfnType = flagSet.String("t", "", "Notification type")
	this.icHostCheckOutput = flagSet.String("o", "", "Check output")
	this.icHostNotfnUser = flagSet.String("b", "", "Manual host notification user")
	this.icHostNotfnComment = flagSet.String("c", "", "Manual host notification comment")
	this.icHostNotfnComment = flagSet.String("i", "", "URL of Webinterface")}

func (this *icingaHostCmd) Parse() Message {
	var msg Message

	msg.Fields = []Field {
		{
			Header: "IPv4",
			Text: *this.icHostV4,
		},
	}
	// TODO

	return msg
}