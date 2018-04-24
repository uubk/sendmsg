package main

type Field struct {
	Text   string
	Header string
	Short  bool
}

type Message struct {
	Head            string
	Body_title      string
	Body_link       string
	Body            string
	Color           string
	Fields          []Field
	Frontend        string
	FrontendIconURL string
}

type Config struct {
	Backend string `yaml:"backend"`
	Webhook string `yaml:"webhook"`
}

type CmdData struct{}
