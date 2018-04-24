package main

import "flag"

type simpleCmdData struct {
	SimpleCFGLocation *string
	head              *string
	title             *string
	title_url         *string
	body              *string
	color             *string
	fields            fieldList
}

func (this *simpleCmdData) Init(flagSet *flag.FlagSet) {
	this.SimpleCFGLocation = flagSet.String("cfg", "/etc/sendmsg.yml", "Path to sendmsg config")
	this.head = flagSet.String("head", "", "The header of the message to send (required)")
	this.title = flagSet.String("title", "", "The title of the message to send (required)")
	this.title_url = flagSet.String("title_url", "", "The url of the title of the message to send")
	this.body = flagSet.String("body", "", "The body of the message to send")
	this.color = flagSet.String("color", "", "The color of the message to send")
	flagSet.Var(&this.fields, "fields", "A comma seperated list of fields (name:text) to be added")
}

func (this *simpleCmdData) Parse() Message {
	var msg Message
	msg.Body = *this.body
	msg.Head = *this.head
	msg.Color = *this.color
	msg.Body_title = *this.title
	msg.Body_link = *this.title_url
	msg.Fields = this.fields

	msg.Frontend = "simple"

	return msg
}
