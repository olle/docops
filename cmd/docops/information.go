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

func (i documentInformation) Catalog() (catalog, error) {
	return collectDocuments(i.root, i.finder, i.parser)
}

func (i documentInformation) Message() (message, error) {
	documents, err := i.Catalog()
	return newMessage(i.name, documents.Titles(kindAsyncAPI)), err
}

type fixedInformation struct {
	msg message
}

func (i fixedInformation) Message() (message, error) {
	return i.msg, nil
}
