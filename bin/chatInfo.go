package main

type ChatInfo struct {
	chat_id           string
	to_id             string
	to_name           string
	to_avatarPath     string
	last_message      string
	last_message_time string
	last_message_id   string
}

func (chat ChatInfo) getString() string {
	return "chatList:" +
		chat.chat_id + "|" +
		chat.to_id + "|" +
		chat.to_name + "|" +
		chat.to_avatarPath + "|" +
		chat.last_message + "|" +
		chat.last_message_time + "|" +
		chat.last_message_id
}
