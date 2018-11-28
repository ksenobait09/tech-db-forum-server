package user

import (
	"github.com/jackc/pgx"
	"tech-db-server/app/database"
	"tech-db-server/app/models/service"
)

var db *pgx.ConnPool

func init() {
	db = database.GetInstance()
}

//easyjson:json
type UserPointList []*User

//easyjson:json
type User struct {
	// Описание пользователя.
	About string `json:"about,omitempty"`

	// Почтовый адрес пользователя (уникальное поле).
	// Required: true
	// Format: email
	Email string `json:"email"`

	// Полное имя пользователя.
	// Required: true
	Fullname string `json:"fullname"`

	// Имя пользователя (уникальное поле).
	// Данное поле допускает только латиницу, цифры и знак подчеркивания.
	// Сравнение имени регистронезависимо.
	//
	// Read Only: true
	Nickname string `json:"nickname,omitempty"`
}

//easyjson:json
type UserUpdate struct {
	// Описание пользователя.
	About interface{} `json:"about,omitempty"`

	// Почтовый адрес пользователя (уникальное поле).
	// Required: true
	// Format: email
	Email interface{} `json:"email"`

	// Полное имя пользователя.
	// Required: true
	Fullname interface{} `json:"fullname"`
}

const sqlInsert = `
	INSERT INTO users (about, email, fullname, nickname)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT DO NOTHING`

const sqlGetByEmailAndNickname = `
	SELECT about, email, fullname, nickname FROM users
	WHERE email = $1
	OR nickname = $2`

const sqlGetByNickname = `
	SELECT about, email, fullname, nickname FROM users
	WHERE nickname = $1`

const sqlUpdate = `
	UPDATE users
	SET about = COALESCE($1, about), 
	email = COALESCE($2, email),
	fullname = COALESCE($3, fullname)
	WHERE nickname = $4`

type Status int

const (
	StatusConflict Status = iota + 1
	StatusNotExist
	StatusOk
)

func (u *User) Create() (user *User, existingUsers UserPointList) {
	existingUsers = nil
	user = nil
	res, _ := db.Exec(sqlInsert, &u.About, &u.Email, &u.Fullname, &u.Nickname)
	if rows := res.RowsAffected(); rows == 0 {
		rows, _ := db.Query(sqlGetByEmailAndNickname, u.Email, u.Nickname)
		existingUsers = make(UserPointList, 0, 1)
		for rows.Next() {
			user := &User{}
			_ = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
			existingUsers = append(existingUsers, user)
		}
		rows.Close()
		return
	}
	user = u
	service.IncUsersCount(1)
	return
}

func Get(nickname string) *User {
	rows, _ := db.Query(sqlGetByNickname, nickname)
	defer rows.Close()
	if rows.Next() {
		user := &User{}
		_ = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		return user
	}
	return nil
}

func Update(nickname string, updateData *UserUpdate) (*User, Status) {
	res, err := db.Exec(sqlUpdate, &updateData.About, &updateData.Email, &updateData.Fullname, nickname)
	if err != nil {
		return nil, StatusConflict
	}
	rows:= res.RowsAffected()
	if rows == 0 {
		return nil, StatusNotExist
	}
	return Get(nickname), StatusOk
}
