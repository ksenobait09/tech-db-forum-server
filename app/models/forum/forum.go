package forum

import (
	"database/sql"
	"tech-db-server/app/database"
	"tech-db-server/app/singletoneLogger"
)

var db *sql.DB

func init() {
	db = database.GetInstance()
}

//easyjson:json
type Forum struct {

	// Общее кол-во сообщений в данном форуме.
	//
	// Read Only: true
	Posts int64 `json:"posts,omitempty"`

	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL), уникальное поле.
	// Required: true
	// Pattern: ^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$
	Slug string `json:"slug"`

	// Общее кол-во ветвей обсуждения в данном форуме.
	//
	// Read Only: true
	Threads int32 `json:"threads,omitempty"`

	// Название форума.
	// Required: true
	Title string `json:"title"`

	// Nickname пользователя, который отвечает за форум.
	// Required: true
	User string `json:"user"`
}

const sqlInsert = `
	INSERT INTO forums (slug, title, "user")
	SELECT forum_data.slug, forum_data.title, COALESCE("users".nickname, forum_data."user")
	FROM (SELECT $1 AS slug, $2 as title, $3 as "user") as forum_data
	LEFT JOIN "users" ON "users".nickname = $4
	ON CONFLICT DO NOTHING
	RETURNING ("user")`

const sqlGetBySlug = `
	SELECT slug, "user", title, threads, posts FROM forums
	WHERE slug = $1 `

type Status int
const (
	StatusConflict Status = iota + 1
	StatusUserNotExist
	StatusOk
)

//easyjson:json
func (forum *Forum) Create() (user *Forum, status Status) {
	err := db.QueryRow(sqlInsert, &forum.Slug,  &forum.Title , &forum.User, &forum.User).Scan(&forum.User)
	if err == sql.ErrNoRows {
		existedForum := Get(forum.Slug)
		return existedForum, StatusConflict
	}
	if err != nil {
		return nil, StatusUserNotExist
	}
	return forum, StatusOk
}

func Get(slug string) *Forum {
	rows, err := db.Query(sqlGetBySlug, slug)
	defer rows.Close()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	if rows.Next() {
		forum := &Forum{}
			err = rows.Scan(&forum.Slug, &forum.User, &forum.Title, &forum.Threads, &forum.Posts)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		return forum
	}
	return nil
}