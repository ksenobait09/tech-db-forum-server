package service

import (
	"database/sql"
	"sync/atomic"
	"tech-db-server/app/database"
	"tech-db-server/app/singletoneLogger"
)

var db *sql.DB

// atomic
var ForumsCount *int32
var PostsCount *int32
var ThreadsCount *int32
var UsersCount *int32

func init() {
	ForumsCount = new(int32)
	PostsCount = new(int32)
	ThreadsCount = new(int32)
	UsersCount = new(int32)
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
	err := db.QueryRow(sqlCounts).Scan(UsersCount, ThreadsCount, ForumsCount, PostsCount)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
}

func GetStatus() *Status {
	currentStatus := &Status{
		Thread: atomic.LoadInt32(ThreadsCount),
		Post: atomic.LoadInt32(PostsCount),
		Forum: atomic.LoadInt32(ForumsCount),
		User: atomic.LoadInt32(UsersCount),
	}
	return currentStatus
}

func IncThreadsCount(increment int) {
	atomic.AddInt32(ThreadsCount, int32(increment))
}

func IncPostsCount(increment int) {
	atomic.AddInt32(PostsCount, int32(increment))
}
func IncForumsCount(increment int) {
	atomic.AddInt32(ForumsCount, int32(increment))
}
func IncUsersCount(increment int) {
	atomic.AddInt32(UsersCount, int32(increment))
}
