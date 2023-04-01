package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// [id] = pipe
var usersOnline map[string]*ConnPipe = make(map[string]*ConnPipe)

type ConnPipe struct {
	con        net.Conn
	id         string
	authorized bool
	replyChan  chan string
}

func newPipe(con net.Conn) {
	pipe := ConnPipe{con, "-1", false, make(chan string)}
	go pipe.handle()
}

func (pipe *ConnPipe) close() {
	pipe.con.Close()
}

func (pipe *ConnPipe) write(s string) {
	pipe.con.Write([]byte(s + "#"))
}

func (pipe *ConnPipe) read() (string, error) {
	buf := make([]byte, 1024)
	rlen, err := pipe.con.Read(buf) // get request
	//error check
	if err != nil {
		pipe.close()
		return "", net.ErrClosed
	}
	return string(buf[:rlen]), nil
}

func (pipe *ConnPipe) handle() {
	defer func() {
		pipe.close()
		delete(usersOnline, pipe.id)
		db.closeSession(pipe.id)
	}()
	pipe.runRequestPipe()
}

func (pipe *ConnPipe) runRequestPipe() {
	for {
		time.Sleep(100 * time.Millisecond)
		read, err := pipe.read()
		if err != nil {
			return
		}
		pipe.parseRequest(read)
	}
}

// Requests may come together
// so we need to split it to single ones
func (pipe *ConnPipe) parseRequest(request string) {

	//requests are separated with '#'
	requestList := strings.Split(request, "#")

	for _, request := range requestList {
		if request == "" {
			continue
		}
		if !pipe.authorized && messageType(request) != "login" {
			pipe.write("NOT AUTHORIZED")
			continue
		}
		fmt.Println("->", request)
		switch messageType(request) {
		///
		case "login":
			reply := login(request)

			if result := messageBody(reply); result != "fail" {
				pipe.authorized = true
				pipe.id = messageBody(reply)
				usersOnline[result] = pipe
			}
			pipe.write(reply)
		///
		case "getChatList":
			for _, chat := range getChats(request) {
				pipe.write(chat.getString())
			}
			///

		case "getMessages":
			for _, message := range getMessages(request) {
				pipe.write(message.getString())
			}
			///
		case "newMessage":
			newMessage(request)
			///
		case "getSessionData":
			pipe.write(getSessionData(request))
			///
		default:
			pipe.write("UNCURRENT REQUEST")
		}
	}
}

func messageType(message string) string {
	s := strings.Split(message, ":")
	if len(s) < 2 {
		return ""
	}
	return s[0]
}
func messageBody(message string) string {
	s := strings.Split(message, ":")
	if len(s) < 2 {
		return ""
	}
	return strings.Join(s[1:], ":")
}
