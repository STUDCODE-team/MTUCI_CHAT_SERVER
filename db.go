package main

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	DB *sql.DB
}

func (db Database) handle() {
	connStr := "admin_root:tTWOKQ0yd0@/admin_mtuci"
	db.DB, _ = sql.Open("mysql", connStr)
	if err := db.DB.Ping(); err != nil {
		log.Panic(err)
	}
}

func (db Database) authUser(login, hashedPassword string) (int, bool) {
	id := 0
	err := db.DB.QueryRow("SELECT id FROM users WHERE login=? AND hash=$?", login, hashedPassword).Scan(&id)
	return id, err == nil
}

type ChatInfo struct {
	chat_id       int
	to_id         int
	to_name       string
	to_avatarPath string
}

func (db Database) getChatList(userID int) []ChatInfo {
	query :=
		`
		SELECT chatList.id, chatList.to_id, userData.name, userData.avatarPath
		FROM chatList 
		JOIN userData
		ON chatList.from_id = ? AND chatList.to_id = userData.id
	`
	rows, _ := db.DB.Query(query, userID)
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
