package main

type MessageInfo struct {
	message string
	fromMe  string
}

func (mess MessageInfo) getString() string {
	return "Message:" +
		mess.message + "|" +
		mess.fromMe
}
