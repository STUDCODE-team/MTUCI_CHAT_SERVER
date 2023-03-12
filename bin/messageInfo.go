package main

type MessageInfo struct {
	message  string
	fromID   string
	fromName string
	avatar   string
	time     string
}

func (mess MessageInfo) getString() string {
	return "Message:" +
		mess.message + "|" +
		mess.fromID + "|" +
		mess.fromName + "|" +
		mess.avatar + "|" +
		mess.time
}
