package main

type message struct {
	greeting       string
	specifications []specificationSummary
}

func newMessage(name string, specifications []specificationSummary) message {
	return message{
		greeting:       "Hello " + name + "!",
		specifications: specifications,
	}
}
