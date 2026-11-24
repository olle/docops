package main

type message struct {
	greeting string
	titles   []string
}

func newMessage(name string, titles []string) message {
	return message{
		greeting: "Hello " + name + "!",
		titles:   titles,
	}
}
