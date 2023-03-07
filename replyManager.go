package main

import (
	"strconv"
	"strings"
)

func login(request string) string {
	requestBoby := messageBody(request)
	login := strings.Split(requestBoby, "|")[0]
	passwordHash := strings.Split(requestBoby, "|")[1]

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
	body := strings.Split(messageBody(request), ":")
	return db.getMessages(body[0], body[1])
}
