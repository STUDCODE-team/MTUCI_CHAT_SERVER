package main

import (
	// "fmt"
	"net"
	// "strconv"
	"strings"
)


func main() {
	go startServer()

	for {

	}
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
		// case <-need_location:
			// con.Write([]byte("GIVEMELOCATION#"))
		// case reply := <-replyChan:
			// location <- reply
		}
	}
}

// Requests may come together
// so we need to split it to single ones
func parseRequest(request string, replyChan chan string) {
	//requests are separated with '#'
	requestList := strings.Split(request, "#")
	for _, singleRequest := range requestList {
		if singleRequest != "" {
			replyChan <- singleRequest
		}
	}
}
