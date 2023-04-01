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
		db.authSession(strconv.Itoa(id))
		return "login:" + strconv.Itoa(id)
	} else {
		return "login:fail"
	}
}

func getChats(request string) []ChatInfo {
	id, _ := strconv.Atoi(messageBody(request))
	return db.getChatList(id)
}

func getSessionData(request string) string {
	id := messageBody(request)
	return "sessionData:" + db.getSessionData(id)
}

func getMessages(request string) []MessageInfo {
	return db.getMessages(messageBody(request))
}

func newMessage(request string) {
	body := strings.Split(messageBody(request), ":")
	chatID := body[0]
	userID := body[1]
	message := strings.Join(body[2:], "")
	m := db.addMessage(chatID, userID, message)
	/// отослать всем пользователям этого чата
	idList := db.getUsersFromChat(chatID)
	// fmt.Println(idList)
	// fmt.Print(usersOnline)
	for _, id := range idList {
		pipe := usersOnline[id]
		if pipe == nil {
			continue
		}
		pipe.write(m.getString())
	}
}
