package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	DB *sql.DB
}

func (db *Database) handle() {
	connStr := "admin_root:tTWOKQ0yd0@tcp(5.253.62.248:3306)/admin_mtuci"
	db.DB, _ = sql.Open("mysql", connStr)
	db.DB.SetMaxIdleConns(0)
	if err := db.DB.Ping(); err != nil {
		log.Panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Second * 15)
			err := db.DB.Ping()
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func (db *Database) authUser(login, hashedPassword string) (int, bool) {
	id := 0
	err := db.DB.QueryRow("SELECT id FROM users WHERE login=? AND hash=?", login, hashedPassword).Scan(&id)
	return id, err == nil
}

func (db *Database) getChatList(userID int) []ChatInfo {
	packet := []ChatInfo{}
	query :=
		`
		SELECT e.id as "CHAT ID", m.user_id as "TO USER", CONCAT(u.name, ' ', u.surname) as "TO NAME", u.avatar as "TO AVATAR",
		mes.message as "TEXT", mes.time as "TIME", mes.id as "MESSAGE ID"
		FROM chats e

		INNER JOIN chats m
		ON e.user_id = ? AND e.id = m.id AND e.user_id != m.user_id AND e.type = 'private'

		JOIN users_data u
  		ON u.id = m.user_id

		JOIN messages mes
		ON mes.id = ( 
			SELECT id
			FROM messages
			WHERE messages.chat_id = e.id
			ORDER BY messages.id DESC
			LIMIT 1
			)
		ORDER BY mes.id ASC
		`
	rows, _ := db.DB.Query(query, userID)
	defer rows.Close()
	for rows.Next() {
		chat := ChatInfo{}
		if err := rows.Scan(&chat.chat_id, &chat.to_id,
			&chat.to_name, &chat.to_avatarPath,
			&chat.last_message, &chat.last_message_time,
			&chat.last_message_id); err != nil {
			log.Fatal(err)
		}
		packet = append(packet, chat)
	}

	query =
		`
		SELECT e.id as "CHAT ID", e.id as "TO USER", u.title as "TITLE", u.avatar as "AVATAR",
		mes.message as "TEXT", mes.time as "TIME", mes.id as "MESSAGE ID"
		FROM chats e

		JOIN chats_data u
  		ON u.id = e.id AND e.user_id = ? AND e.type = 'group'

		JOIN messages mes
		ON mes.id = ( 
			SELECT id
			FROM messages
			WHERE messages.chat_id = e.id
			ORDER BY messages.id DESC
			LIMIT 1
			)
		ORDER BY mes.id ASC
		`
	rows, _ = db.DB.Query(query, userID)

	for rows.Next() {
		chat := ChatInfo{}
		if err := rows.Scan(&chat.chat_id, &chat.to_id,
			&chat.to_name, &chat.to_avatarPath,
			&chat.last_message, &chat.last_message_time,
			&chat.last_message_id); err != nil {
			log.Fatal(err, chat)
		}
		packet = append(packet, chat)

	}
	return packet
}

func (db *Database) getMessages(chatID string) []MessageInfo {
	packet := []MessageInfo{}
	query :=
		`
		SELECT m.message, u.id, CONCAT(u.name, " ", u.surname), u.avatar, m.time
		FROM messages m
		JOIN users_data u
		ON m.chat_id = ? AND m.user_id = u.id
		ORDER BY m.id ASC
		`
	rows, _ := db.DB.Query(query, chatID)
	defer rows.Close()
	for rows.Next() {
		m := MessageInfo{}
		err := rows.Scan(&m.message, &m.fromID,
			&m.fromName, &m.avatar, &m.time)
		if err != nil {
			log.Fatal(err)
		}
		packet = append(packet, m)
	}
	return packet
}

func (db *Database) addMessage(chatID, userID, message string) MessageInfo {
	query :=
		`
		INSERT INTO messages (chat_id, user_id, message) VALUES (?,?,?)
		`
	result, _ := db.DB.Exec(query, chatID, userID, message)
	id, _ := result.LastInsertId()
	query =
		`
		SELECT m.message, u.id, CONCAT(u.name, " ", u.surname), u.avatar, m.time
		FROM messages m
		JOIN users_data u
		ON m.id = ? AND m.user_id = u.id
		`

	row := db.DB.QueryRow(query, id)
	m := MessageInfo{}
	err := row.Scan(&m.message, &m.fromID,
		&m.fromName, &m.avatar, &m.time)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (db *Database) getUsersFromChat(chatID string) []string {
	idList := make([]string, 0)
	query :=
		`
		select user_id from chats where id = ?
	`
	rows, _ := db.DB.Query(query, chatID)

	for rows.Next() {
		var id string = ""
		rows.Scan(&id)
		idList = append(idList, id)
	}
	return idList
}
