package main

type informationSource interface {
	Message() (message, error)
}

type catalogSource interface {
	Catalog() (catalog, error)
}

type documentInformation struct {
	name    string
	root    string
	finder  documentFinder
	sniffer documentSniffer
	parser  documentParser
}

func (i documentInformation) Catalog() (catalog, error) {
	return collectDocuments(i.root, i.finder, i.sniffer, i.parser)
}

func (i documentInformation) Message() (message, error) {
	documents, err := i.Catalog()
	return newMessage(i.name, documents.Summaries()), err
}

type fixedInformation struct {
	msg message
}

func (i fixedInformation) Message() (message, error) {
	return i.msg, nil
}
