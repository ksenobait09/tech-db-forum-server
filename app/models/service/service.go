package service

import (
	"database/sql"
	"tech-db-server/app/database"
	"tech-db-server/app/singletoneLogger"
)

var db *sql.DB

// mutex
var ForumsCount int
var PostsCount int
var ThreadsCount int
var UsersCount int

func init() {
	db = database.GetInstance()
	//initStatus()
}

//easyjson:json
type Status struct {
	// Кол-во разделов в базе данных.
	// Required: true
	Forum int `json:"forum"`

	// Кол-во сообщений в базе данных.
	// Required: true
	Post int `json:"post"`

	// Кол-во веток обсуждения в базе данных.
	// Required: true
	Thread int `json:"thread"`

	// Кол-во пользователей в базе данных.
	// Required: true
	User int `json:"user"`
}

const sqlClear = `
	TRUNCATE users, forums, threads, posts, votes;`

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
	err := db.QueryRow(sqlCounts).Scan(&UsersCount, &ThreadsCount, &ForumsCount, &PostsCount)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
}

func GetStatus() *Status {
	initStatus()
	currentStatus := &Status{}
	currentStatus.Thread = ThreadsCount
	currentStatus.Post = PostsCount
	currentStatus.Forum = ForumsCount
	currentStatus.User = UsersCount
	return currentStatus
}
