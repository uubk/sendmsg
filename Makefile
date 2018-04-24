build:
	go build -i -v -ldflags="-X main.gitversion=$(shell git describe --always --long --dirty)"
	strip sendmsg
