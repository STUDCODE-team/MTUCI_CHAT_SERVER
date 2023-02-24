package main

import (
	// "fmt"
	"fmt"
	"strconv"

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

func authUser(login, hashedPassword string) (int, bool) {
	id := 0
	err := DB.QueryRow("SELECT id FROM users WHERE login=$1 AND password=$2", login, hashedPassword).Scan(&id)
	return id, err == nil
}

type ChatInfo struct {
	chat_id       int
	to_id         int
	to_name       string
	to_avatarPath string
}

func getChatList(userID int) []ChatInfo {
	query :=
		`
		SELECT chatList.id, chatList.to_id, userData.name, userData.avatarPath
		FROM chatList 
		JOIN userData
		ON chatList.from_id = $1 AND chatList.to_id = userData.id
	`
	rows, _ := DB.Query(query, userID)
	defer rows.Close()

	packet := []ChatInfo{}

	for rows.Next() {
		chat := ChatInfo{}

		if err := rows.Scan(&chat.chat_id, &chat.to_id,
			&chat.to_name, &chat.to_avatarPath); err != nil {
			log.Fatal(err)
		}
		packet = append(packet, chat)
	}
	return packet
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

			if id, ok := authUser(login, passwordHash); ok {
				replyChan <- "login:" + strconv.Itoa(id)
			} else {
				replyChan <- "login:fail"
			}
		///
		///
		///
		case "getChatList":
			id, _ := strconv.Atoi(messageBody(request))
			packet := getChatList(id)
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
