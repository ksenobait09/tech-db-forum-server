package forum

import (
	"database/sql"
	"fmt"
	"strings"
	"tech-db-server/app/database"
	"tech-db-server/app/models/user"
	"tech-db-server/app/models/service"
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
	WHERE slug = $1`


const sqlSelectUserForum = `
	SELECT u.nickname, u.email, u.about, u.fullname
	FROM userforum AS uf
	JOIN users AS u ON u.nickname = uf.nickname
	WHERE uf.slug = $1
`

const sqlInsertUserForum = `
INSERT INTO userforum (slug, nickname) VALUES ($1, $2) 
ON CONFLICT DO NOTHING`

type Status int

const (
	StatusConflict Status = iota + 1
	StatusSomethingNotExist
	StatusOk
)

func (forum *Forum) Create() (*Forum, Status) {
	err := db.QueryRow(sqlInsert, &forum.Slug, &forum.Title, &forum.User, &forum.User).Scan(&forum.User)
	if err == sql.ErrNoRows {
		existedForum := Get(forum.Slug)
		return existedForum, StatusConflict
	}
	if err != nil {
		return nil, StatusSomethingNotExist
	}
	service.IncForumsCount(1)
	return forum, StatusOk
}

func Get(slug string) *Forum {
	rows, _ := db.Query(sqlGetBySlug, slug)
	defer rows.Close()

	if rows.Next() {
		forum := &Forum{}
		rows.Scan(&forum.Slug, &forum.User, &forum.Title, &forum.Threads, &forum.Posts)

		return forum
	}
	return nil
}

func GetUsers(slug string, limit int, since string, desc bool) (user.UserPointList, Status) {
	if !IsForumExists(slug) {
		return nil, StatusSomethingNotExist
	}
	users := make(user.UserPointList, 0, limit)
	var query strings.Builder
	query.Grow(100)
	fmt.Fprint(&query, sqlSelectUserForum)
	if since != "" {
		if desc {
			fmt.Fprint(&query, " AND uf.nickname < $2")
		} else {
			fmt.Fprint(&query, " AND uf.nickname > $2 ")
		}
	} else {
		fmt.Fprint(&query, " AND $2 = ''")
	}
	if desc {
		fmt.Fprint(&query, " ORDER BY uf.nickname DESC")
	} else {
		fmt.Fprint(&query, " ORDER BY uf.nickname ASC")
	}
	if limit > 0 {
		fmt.Fprint(&query, " LIMIT $3")
	} else {
		fmt.Fprint(&query, " LIMIT 100000+$3")
	}
	rows, _ := db.Query(query.String(), slug, since, limit)

	for rows.Next() {
		u := &user.User{}
		rows.Scan(&u.Nickname, &u.Email, &u.About, &u.Fullname)
		users = append(users, u)
	}
	rows.Close()
	return users, StatusOk
}

func IsForumExists(slug string) bool {
	err := db.QueryRow(`SELECT slug FROM forums WHERE slug=$1`, slug).Scan(&slug)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

func InsertMapIntoUserForum(tx *sql.Tx, slug string, users map[string]bool) {
	lenUsers := len(users)
	if lenUsers == 0 {
		return
	}
	var query strings.Builder
	query.Grow(46 + 11*lenUsers + 23)
	fmt.Fprint(&query, "INSERT INTO userforum (slug, nickname) VALUES")
	counterPlaceholders := 1
	i := 1
	var args []interface{}
	for u := range users {
		first := counterPlaceholders
		counterPlaceholders++
		second := counterPlaceholders
		counterPlaceholders++
		if lenUsers == i {
			fmt.Fprintf(&query, " ($%d, $%d)", first, second)
		} else {
			fmt.Fprintf(&query, " ($%d, $%d),", first, second)
		}
		args = append(args, slug, u)
		i++
	}
	fmt.Fprintf(&query, " ON CONFLICT DO NOTHING")
	_, err := tx.Exec(query.String(), args...)
	if err != nil {
		tx.Rollback()
	}
}

func InsertIntoUserForum(tx *sql.Tx, slug string, user string) {
	_, err := tx.Exec(sqlInsertUserForum, slug, user)
	if err != nil {
		tx.Rollback()
	}
}
