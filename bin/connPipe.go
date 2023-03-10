package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type ConnPipe struct {
	con        net.Conn
	authorized bool
	replyChan  chan string
}

func newPipe(con net.Conn) {
	pipe := ConnPipe{con, false, make(chan string)}
	go pipe.handle()
}

func (pipe *ConnPipe) close() {
	pipe.con.Close()
}

func (pipe *ConnPipe) write(s string) {
	pipe.con.Write([]byte(s + "#"))
}

func (pipe *ConnPipe) read() string {
	buf := make([]byte, 1024)
	rlen, err := pipe.con.Read(buf) // get request
	//error check
	if err != nil {
		return ""
	}
	return string(buf[:rlen])
}

func (pipe *ConnPipe) handle() {
	defer pipe.close()
	pipe.runRequestPipe()
}

func (pipe *ConnPipe) runRequestPipe() {
	for {
		time.Sleep(100 * time.Millisecond)
		pipe.parseRequest(pipe.read())
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
			pipe.authorized = (messageBody(reply) != "fail")
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
