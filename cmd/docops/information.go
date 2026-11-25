package main

type informationSource interface {
	Message() (message, error)
}

type documentInformation struct {
	name   string
	root   string
	finder documentFinder
	parser documentParser
}

func (i documentInformation) Message() (message, error) {
	titles, err := collectTitles(i.root, i.finder, i.parser)
	return newMessage(i.name, titles), err
}

type fixedInformation struct {
	msg message
}

func (i fixedInformation) Message() (message, error) {
	return i.msg, nil
}
