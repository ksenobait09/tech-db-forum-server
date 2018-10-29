package service

import (
	"database/sql"
	"tech-db-server/app/database"
	"tech-db-server/app/singletoneLogger"
)

var db *sql.DB
var currentStatus *Status

var ForumsCount *int32
var PostsCount *int32
var ThreadsCount *int32
var UsersCount *int32

func init() {
	currentStatus = &Status{}
	db = database.GetInstance()
	initStatus()
}

//easyjson:json
type Status struct {
	// Кол-во разделов в базе данных.
	// Required: true
	Forum int32 `json:"forum"`

	// Кол-во сообщений в базе данных.
	// Required: true
	Post int32 `json:"post"`

	// Кол-во веток обсуждения в базе данных.
	// Required: true
	Thread int32 `json:"thread"`

	// Кол-во пользователей в базе данных.
	// Required: true
	User int32 `json:"user"`
}

const sqlClear = `
	TRUNCATE users;
	TRUNCATE forums;
	TRUNCATE threads;
	TRUNCATE posts;
	TRUNCATE voices;`

const sqlCounts = `
SELECT *
FROM (SELECT COUNT(*) FROM "users") as "users"
CROSS JOIN (SELECT COUNT(*) FROM "threads") as threads
CROSS JOIN (SELECT COUNT(*) FROM "forums") as forums
CROSS JOIN (SELECT COUNT(*) FROM "posts") as posts`

func ClearDatabase() {
	_, err := db.Exec(sqlClear)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
}

func initStatus() {
	err := db.QueryRow(sqlCounts).Scan(ForumsCount, PostsCount,	ThreadsCount, UsersCount)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
}

func GetStatus() (*Status) {
	initStatus()
	currentStatus.Thread = *ThreadsCount
	currentStatus.Post = *PostsCount
	currentStatus.Forum = *ForumsCount
	currentStatus.User = *UsersCount
	return currentStatus
}