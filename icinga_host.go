package main

import (
	"flag"
	"bytes"
	"net/url"
	"time"
	"github.com/sirupsen/logrus"
	"strconv"
)

var icGoodState = map[string]bool {
	"UP": true,
}

var icBadState = map[string]bool {
	"UNREACHABLE": true,
	"DOWN": true,
}

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
	// Fields below are unused currently but need to be parsed in order to be compliant
	icDateTime *string
	icHostDisplayname *string
	icUserEmail *string
	icMailFrom *string
	icToSyslog *string
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
	this.icURL = flagSet.String("i", "", "URL of Webinterface")
	this.icDateTime = flagSet.String("d", "", "Date/Time of event")
	this.icHostDisplayname = flagSet.String("n", "", "Host display name")
	this.icUserEmail = flagSet.String("r", "", "User email")
	this.icMailFrom = flagSet.String("f", "", "Sender email")
	this.icToSyslog = flagSet.String("v", "", "Send to syslog")
}

func (this *icingaHostCmd) Parse() Message {
	var msg Message

	msg.Body_title = *this.icHostname + " is " + *this.icHostState

	addFieldToMessage(&msg, true, "IPv4", this.icHostV4)
	addFieldToMessage(&msg, true, "IPv6", this.icHostV6)
	addFieldToMessage(&msg, true, "Notification type", this.icHostNotfnType)
	addFieldToMessage(&msg, true, "Notification user", this.icHostNotfnUser)
	timestamp, err := time.Parse("2006-01-02 15:04:05 -0700", *this.icDateTime)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Couldn't parse supplied timestamp")
	} else {
		timeStr := convertToSlackDate(timestamp)
		addFieldToMessage(&msg, true, "Timestamp", &timeStr)
	}

	var buffer bytes.Buffer

	buffer.WriteString("```\n")
	buffer.WriteString(*this.icHostCheckOutput)
	buffer.WriteString("```\n")
	if *this.icHostNotfnComment != "" {
		buffer.WriteString("Comment: ")
		buffer.WriteString(*this.icHostNotfnComment)
	}
	msg.Body = buffer.String()

	if *this.icURL != "" {
		params := url.Values{}
		params.Add("host", *this.icHostname)
		msg.Body_link = *this.icURL + "/monitoring/host/show?"+ params.Encode()
	}

	if icBadState[*this.icHostState] {
		msg.Color = "#ff5566"
	} else if icGoodState[*this.icHostState] {
		msg.Color = "#44bb77"
	} else {
		msg.Color = "#aa44ff"
	}

	msg.Frontend = "Icinga2 Host Notification"
	msg.FrontendIconURL = "https://raw.githubusercontent.com/Icinga/icingaweb2/master/public/img/favicon.png"

	return msg
}

func addFieldToMessage(msg *Message, short bool, header string, field *string) {
	if *field != "" {
		msg.Fields = append(msg.Fields, Field{
			Header: header,
			Text: *field,
			Short: short,
		})
	}
}

func convertToSlackDate(timestamp time.Time) string {
	return "<!date^" + strconv.FormatInt(timestamp.Unix(), 10) + "^{date_num} {time_secs}|" +
		timestamp.Format("Mon Jan 2 15:04:05 -0700 MST 2006")+ ">"
}