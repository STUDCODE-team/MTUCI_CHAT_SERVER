package main

import (
	"net"
	"strings"
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

	go pipe.runRequestPipe()
	pipe.runReplyPipe()
}

func (pipe *ConnPipe) runRequestPipe() {
	for {
		go pipe.parseRequest(pipe.read())
	}
}

func (pipe *ConnPipe) runReplyPipe() {
	for {
		select {
		case reply := <-pipe.replyChan:

			switch messageType(reply) {

			case "login":
				pipe.authorized = (messageBody(reply) != "fail")
				pipe.write(reply)

			case "chatList":
				pipe.write(reply)

			}
		}
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
			continue
		}

		switch messageType(request) {
		///
		case "login":
			pipe.replyChan <- login(request)
		///

		case "getChatList":
			for _, chat := range getChats(request) {

				pipe.replyChan <- chat.getString()
			}
		}
		///
	}
}

func messageType(message string) string {
	return strings.Split(message, ":")[0]
}
func messageBody(message string) string {
	return strings.Join(strings.Split(message, ":")[1:], "")
}
