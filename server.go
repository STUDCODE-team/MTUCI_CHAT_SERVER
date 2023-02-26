package main

import (
	// "fmt"
	"fmt"
	"strconv"

	// "log"
	"net"

	// "strconv"
	"strings"
)

var db Database

func main() {
	db.handle()
	startServer()
}

func startServer() {
	// server creation
	dstream, err := net.Listen("tcp", ":30391")
	if err != nil {
		return
	}
	defer dstream.Close()

	// handle new connections in a loop
	for {
		// accept new connection
		con, err := dstream.Accept()
		if err != nil {
			return
		}
		// procced connection above in separated virtual thread
		go handle(con)
	}
}

func handle(con net.Conn) {
	defer con.Close()
	// create new channel to send replies
	replyChan := make(chan string)

	//points if user has authorized
	userAuthorized := false

	// get new client requests in loop in new thread
	go func() {
		for {
			buf := make([]byte, 128)
			rlen, err := con.Read(buf) // get request
			//error check
			if err != nil {
				return
			}
			// send request pack to parse it via function
			go parseRequest(string(buf[:rlen]), replyChan)
		}
	}()

	//sending replies to client in the loop
	for {
		select {
		case reply := <-replyChan:

			if !userAuthorized && messageType(reply) != "login" {
				con.Write([]byte("login:fail"))
				continue
			}

			switch messageType(reply) {

			case "login":
				userAuthorized = (messageBody(reply) != "fail")
				con.Write([]byte(reply))

			case "chatList":
				fmt.Println(reply)
				con.Write([]byte(reply))

			}

		}
	}
}

// Requests may come together
// so we need to split it to single ones
func parseRequest(request string, replyChan chan string) {

	//requests are separated with '#'
	requestList := strings.Split(request, "#")

	for _, request := range requestList {

		switch messageType(request) {

		///
		case "login":
			requestBoby := messageBody(request)
			login := strings.Split(requestBoby, "|")[0]
			passwordHash := strings.Split(requestBoby, "|")[1]

			if id, ok := db.authUser(login, passwordHash); ok {
				replyChan <- "login:" + strconv.Itoa(id)
			} else {
				replyChan <- "login:fail"
			}
		///
		///
		///
		case "getChatList":
			id, _ := strconv.Atoi(messageBody(request))
			packet := db.getChatList(id)
			for _, item := range packet {
				replyChan <- "chatList:" +
					strconv.Itoa(item.chat_id) + ":" +
					strconv.Itoa(item.to_id) + ":" +
					item.to_name + ":" +
					item.to_avatarPath
			}
		}
		///
		///
		///
	}
}

func messageType(message string) string {
	return strings.Split(message, ":")[0]
}

func messageBody(message string) string {
	return strings.Split(message, ":")[1]
}
