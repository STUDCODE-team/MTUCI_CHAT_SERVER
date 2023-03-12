package main

import (
	"strconv"
	"strings"
)

func login(request string) string {
	requestBoby := messageBody(request)
	s := strings.Split(requestBoby, "|")
	if len(s) < 2 {
		return "login:fail"
	}
	login := s[0]
	passwordHash := s[1]

	if id, ok := db.authUser(login, passwordHash); ok {
		return "login:" + strconv.Itoa(id)
	} else {
		return "login:fail"
	}
}

func getChats(request string) []ChatInfo {
	id, _ := strconv.Atoi(messageBody(request))
	return db.getChatList(id)
}

func getMessages(request string) []MessageInfo {
	return db.getMessages(messageBody(request))
}

func newMessage(request string) string {
	body := strings.Split(messageBody(request), ":")
	m := db.addMessage(body[0], body[1], body[2])
	return m.getString()
}
