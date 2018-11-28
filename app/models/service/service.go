package service

import (
	"github.com/jackc/pgx"
	"sync/atomic"
	"tech-db-server/app/database"
)

var db *pgx.ConnPool

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
	TRUNCATE users, forums, threads, posts, votes, userforum;`

const sqlCounts = `
SELECT *
FROM (SELECT COUNT(*) FROM "users") as "users"
CROSS JOIN (SELECT COUNT(*) FROM "threads") as threads
CROSS JOIN (SELECT COUNT(*) FROM "forums") as forums
CROSS JOIN (SELECT COUNT(*) FROM "posts") as posts`

func ClearDatabase() {
	 db.Exec(sqlClear)
	resetThreadsCount()
	resetPostsCount()
	resetForumsCount()
	resetUsersCount()
}

func initStatus() {
	db.QueryRow(sqlCounts).Scan(UsersCount, ThreadsCount, ForumsCount, PostsCount)
}

func getStatus() *Status {
	status := &Status{}
	 db.QueryRow(sqlCounts).Scan(&status.User, &status.Thread, &status.Forum, &status.Post)
	return status
}
func GetStatus() *Status {
	currentStatus := &Status{
		Thread: atomic.LoadInt32(ThreadsCount),
		Post:   atomic.LoadInt32(PostsCount),
		Forum:  atomic.LoadInt32(ForumsCount),
		User:   atomic.LoadInt32(UsersCount),
	}
	return currentStatus
	//return getStatus()
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

func resetThreadsCount() {
	atomic.StoreInt32(ThreadsCount, 0)
}
func resetPostsCount() {
	atomic.StoreInt32(PostsCount, 0)
}
func resetForumsCount() {
	atomic.StoreInt32(ForumsCount, 0)
}
func resetUsersCount() {
	atomic.StoreInt32(UsersCount, 0)
}
