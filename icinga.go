package main

import (
	"flag"
	"bytes"
	"net/url"
	"time"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var icGoodState = map[string]bool {
	"UP": true,
	"OK": true,
}

var icBadState = map[string]bool {
	"UNREACHABLE": true,
	"DOWN": true,
	"CRITICAL": true,
}

var icWarnState = map[string]bool {
	"WARNING": true,
}

type icingaCmd struct {
	SimpleCFGLocation *string
	// Shared fields
	icHostname     *string
	icHostV4       *string
	icHostV6       *string
	icCheckState   *string
	icCheckOutput  *string
	icNotfnType    *string
	icNotfnUser    *string
	icNotfnComment *string
	icURL          *string
	// Service-only fields
	icServiceName *string
	icServiceDisplayName *string
	// Fields below are unused currently but need to be parsed in order to be compliant
	icDateTime *string
	icHostDisplayname *string
	icUserEmail *string
	icMailFrom *string
	icToSyslog *string
	// Which command are we executing?
	icIsServiceCMD bool
}

func (this *icingaCmd) Init(flagSet *flag.FlagSet, isServiceCMD bool) {
	this.SimpleCFGLocation = flagSet.String("cfg", "/etc/sendmsg.yml", "Path to sendmsg config")
	this.icHostname = flagSet.String("l", "", "Hostname")
	this.icHostV4 = flagSet.String("4", "", "Host address (v4)")
	this.icHostV6 = flagSet.String("6", "", "Host address (v6)")
	this.icCheckState = flagSet.String("s", "", "Host state")
	this.icNotfnType = flagSet.String("t", "", "Notification type")
	this.icCheckOutput = flagSet.String("o", "", "Check output")
	this.icNotfnUser = flagSet.String("b", "", "Manual notification user")
	this.icNotfnComment = flagSet.String("c", "", "Manual notification comment")
	this.icURL = flagSet.String("i", "", "URL of Webinterface")
	this.icDateTime = flagSet.String("d", "", "Date/Time of event")
	this.icHostDisplayname = flagSet.String("n", "", "Host display name")
	this.icUserEmail = flagSet.String("r", "", "User email")
	this.icMailFrom = flagSet.String("f", "", "Sender email")
	this.icToSyslog = flagSet.String("v", "", "Send to syslog")
	this.icIsServiceCMD = isServiceCMD
	if isServiceCMD {
		this.icServiceName = flagSet.String("e", "", "Service name")
		this.icServiceDisplayName = flagSet.String("u", "", "Service display name")
	}
}

func (this *icingaCmd) Parse() Message {
	var msg Message

	if this.icIsServiceCMD {
		msg.Body_title = *this.icServiceName + " is " + *this.icCheckState
	} else {
		msg.Body_title = *this.icHostname + " is " + *this.icCheckState
	}

	addFieldToMessage(&msg, true, "IPv4", this.icHostV4)
	addFieldToMessage(&msg, true, "IPv6", this.icHostV6)
	addFieldToMessage(&msg, true, "Notification type", this.icNotfnType)
	addFieldToMessage(&msg, true, "Notification user", this.icNotfnUser)
	if this.icIsServiceCMD {
		addFieldToMessage(&msg, true, "Host", this.icHostname)
	}
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
	buffer.WriteString(*this.icCheckOutput)
	buffer.WriteString("```\n")
	if *this.icNotfnComment != "" {
		buffer.WriteString("Comment: ")
		buffer.WriteString(*this.icNotfnComment)
	}
	msg.Body = buffer.String()

	if *this.icURL != "" {
		params := url.Values{}
		params.Add("host", *this.icHostname)
		if this.icIsServiceCMD {
			params.Add("service", *this.icServiceName)
			msg.Body_link = *this.icURL + "/monitoring/service/show?"+ strings.Replace(params.Encode(), "+", "%20", -1)
		} else {
			msg.Body_link = *this.icURL + "/monitoring/host/show?"+ strings.Replace(params.Encode(), "+", "%20", -1)
		}
	}

	if icBadState[*this.icCheckState] {
		msg.Color = "#ff5566"
	} else if icGoodState[*this.icCheckState] {
		msg.Color = "#44bb77"
	} else if icWarnState[*this.icCheckState] {
		msg.Color = "#ffaa44"
	} else {
		msg.Color = "#aa44ff"
	}

	if this.icIsServiceCMD {
		msg.Frontend = "Icinga2 Service Notification"
	} else {
		msg.Frontend = "Icinga2 Host Notification"
	}
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