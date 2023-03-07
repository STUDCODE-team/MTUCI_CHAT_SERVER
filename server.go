package main

import (
	"net"
	// "strconv"
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
		newPipe(con)
	}
}
