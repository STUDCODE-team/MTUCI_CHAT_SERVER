package main

import (
	// "fmt"

	"log"
	"net"

	// "strconv"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func connectDB() {
	connStr := "user=toniess password=123456 dbname=mtuci_chat sslmode=disable"
	DB, _ = sql.Open("postgres", connStr)
	if err := DB.Ping(); err != nil {
		log.Panic(err)
	}
}

func authUser(login, hashedPassword string) bool {
	id := 0
	err := DB.QueryRow("SELECT id FROM users WHERE login=$1 AND password=$2", login, hashedPassword).Scan(&id)
	return err == nil
}

func main() {
	connectDB()
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

	if userAuthorized {
		log.Fatal()
	}

	//sending replies to client in the loop
	for {
		select {
		case reply := <-replyChan:
			switch messageType(reply) {

			case "login":
				userAuthorized = (messageBody(reply) == "ok")
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
		///
		///
		///
		case "login":
			requestBoby := messageBody(request)
			login := strings.Split(requestBoby, "|")[0]
			passwordHash := strings.Split(requestBoby, "|")[1]
			if authUser(login, passwordHash) {
				replyChan <- "login:ok"
			} else {
				replyChan <- "login:fail"
			}

		}
	}
}

func messageType(message string) string {
	return strings.Split(message, ":")[0]
}

func messageBody(message string) string {
	return strings.Split(message, ":")[1]
}
