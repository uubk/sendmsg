package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SlackField struct {
	Text  string `json:"value,omitempty"`
	Title string `json:"title,omitempty"`
	Short bool   `json:"short,omitempty"`
}

type SlackAttachments struct {
	Color      string       `json:"color,omitempty"`
	Text       string       `json:"text,omitempty"`
	Title      string       `json:"title,omitempty"`
	TitleLink  string       `json:"title_link,omitempty"`
	Fields     []SlackField `json:"fields,omitempty"`
	Footer     string       `json:"footer,omitempty"`
	FooterIcon string       `json:"footer_icon,omitempty"`
	Head       string       `json:"author_name,omitempty"`
}

type SlackMessage struct {
	Text        string             `json:"text,omitempty"`
	Attachments []SlackAttachments `json:"attachments,omitempty"`
}

func send_with_slack(msg Message, cfg Config) {
	var body SlackMessage
	var attachement SlackAttachments
	attachement.Head = msg.Head
	attachement.Text = msg.Body
	if msg.Color != "" {
		attachement.Color = msg.Color
	}
	if msg.Body_title != "" {
		attachement.Title = msg.Body_title
	}
	if msg.Body_link != "" {
		attachement.TitleLink = msg.Body_link
	}
	var fields []SlackField
	for _, value := range msg.Fields {
		var field SlackField
		field.Title = value.Header
		field.Text = value.Text
		field.Short = value.Short
		fields = append(fields, field)
	}
	attachement.Fields = fields
	attachement.Footer = msg.Frontend + " (sendmsg "+ gitversion+ ")"
	if msg.FrontendIconURL != "" {
		attachement.FooterIcon = msg.FrontendIconURL
	}
	body.Attachments = []SlackAttachments{attachement}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("Coudln't masrshal message for JSON transport: ", err)
		return
	}
	resp, err := http.Post(cfg.Webhook, "application/reader", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Errorf("Coudln't POST message to webhook: ", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Couldn't read body", err)
	}
	bodyString := string(bodyBytes)
	fmt.Print(bodyString)

	if bodyString != "ok" {
		fmt.Errorf("Slack didn't like what we sent, error: ", bodyString)
	}
}
