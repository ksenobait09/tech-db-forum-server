package thread

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	"tech-db-server/app/database"
	"tech-db-server/app/models/forum"
	"tech-db-server/app/singletoneLogger"
	"tech-db-server/app/models/service"
)

var db *sql.DB

func init() {
	db = database.GetInstance()
}

//easyjson:json
type ThreadPointList []*Thread

//easyjson:json
type Thread struct {
	// Пользователь, создавший данную тему.
	// Required: true
	Author string `json:"author"`

	// Дата создания ветки на форуме.
	// Format: date-time
	Created *strfmt.DateTime `json:"created,omitempty"`

	// Форум, в котором расположена данная ветка обсуждения.
	// Read Only: true
	Forum string `json:"forum,omitempty"`

	// Идентификатор ветки обсуждения.
	// Read Only: true
	ID int `json:"id,omitempty"`

	// Описание ветки обсуждения.
	// Required: true
	Message string `json:"message"`

	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL).
	// В данной структуре slug опционален и не может быть числом.
	//
	// Read Only: true
	// Pattern: ^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$
	Slug string `json:"slug,omitempty"`

	// Заголовок ветки обсуждения.
	// Required: true
	Title string `json:"title"`

	// Кол-во голосов непосредственно за данное сообщение форума.
	// Read Only: true
	Votes int32 `json:"votes,omitempty"`
}

//easyjson:json
type ThreadUpdate struct {

	// Описание ветки обсуждения.
	Message interface{} `json:"message,omitempty"`

	// Заголовок ветки обсуждения.
	Title interface{} `json:"title,omitempty"`
}

//easyjson:json
type Vote struct {
	// Идентификатор пользователя.
	// Required: true
	Nickname string `json:"nickname"`

	// Отданный голос.
	// Required: true
	// Enum: [-1 1]
	Voice int8 `json:"voice"`
}

const sqlInsert = `
	INSERT INTO threads (author, created, forum, message, slug, title)
	SELECT COALESCE("users".nickname, t.author), t.created, COALESCE("forums".slug, t.forum), t.message, t.slug, t.title
	FROM (SELECT $1 as author, $2::timestamptz as created, $3 as forum, $4 as message, $5 as slug, $6 as title) as t
	LEFT JOIN "users" ON "users".nickname = $7
	LEFT JOIN forums ON forums.slug = $8
	ON CONFLICT DO NOTHING
	RETURNING author, created, forum, message, slug, title, id `

const sqlGetBySlug = `
	SELECT author, created, forum, "message" , slug, title, id, votes
	FROM threads
	WHERE slug = $1 
`

const sqlGetById = `
	SELECT author, created, forum, "message", slug, title, id, votes
	FROM threads
	WHERE id = $1 
`

const sqlUpdateById = `
	UPDATE threads
	SET "message" = COALESCE($1, "message"),
	title = COALESCE($2, title)
	WHERE id = $3
	RETURNING author, created, forum, "message", slug, title, id, votes
`
const sqlUpdateBySlug = `
	UPDATE threads
	SET "message" = COALESCE($1, "message"),
	title = COALESCE($2, title)
	WHERE slug = $3
	RETURNING author, created, forum, "message" , slug, title, id, votes
`

const sqlGetByForumSlug = `
	SELECT author, created, forum, "message", slug, title, id, votes
	FROM threads
	WHERE forum = $1
`

const sqlInsertVote = `
	INSERT INTO votes (idThread, nickname, voice) VALUES ($1, $2, $3)
`

const sqlUpdateVote = `
	UPDATE votes SET 
	voice = $3
	WHERE idThread = $1 
	AND nickname = $2
`

const sqlUpdateThreadVotes = `
	UPDATE threads SET
	votes = $1
	WHERE id = $2
	RETURNING author, created, forum, "message" , slug, title, id, votes
`

const sqlSelectThreadAndVoteBySlug = `
	SELECT votes.voice, threads.id, threads.votes, u.nickname
	FROM (SELECT 1) s
	LEFT JOIN threads ON threads.slug = $1
	LEFT JOIN users as u ON u.nickname = $2
	LEFT JOIN votes ON threads.id = votes.iDthread AND u.nickname = votes.nickname
`

const sqlSelectThreadAndVoteById = `
	SELECT votes.voice, threads.id, threads.votes, u.nickname
	FROM (SELECT 1) s
	LEFT JOIN threads ON threads.id = $1
	LEFT JOIN "users" u ON u.nickname = $2
	LEFT JOIN votes ON threads.id = votes.iDthread AND u.nickname = votes.nickname
`

const sqlUpdateThreadsCount = `
	UPDATE forums
	SET threads = threads + 1
	WHERE slug = $1`

const sqlGetThreadIdBySlug = `
	SELECT id
	FROM threads
	WHERE slug = $1
`
const sqlGetThreadIdById = `
	SELECT id
	FROM threads
	WHERE id = $1
`

type Status int

const (
	StatusConflict Status = iota + 1
	StatusUserOrForumNotExist
	StatusOk
	StatusNotFound
)

func createSqlNullString(str string) *sql.NullString {
	return &sql.NullString{String: str, Valid: str != ""}
}

func (thread *Thread) Create() Status {
	tx, err := db.Begin()
	defer tx.Rollback()
	slugNullable := createSqlNullString(thread.Slug)
	err = tx.QueryRow(sqlInsert, &thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.Author, &thread.Forum).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID)
	thread.Slug = slugNullable.String
	if err == sql.ErrNoRows {
		thread.Get(thread.Slug, 0)
		return StatusConflict
	}
	if err != nil {
		return StatusUserOrForumNotExist
	}
	_, err = tx.Exec(sqlUpdateThreadsCount, thread.Forum)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusUserOrForumNotExist
	}
	forum.InsertIntoUserForum(tx, thread.Forum, thread.Author)
	tx.Commit()
	service.IncThreadsCount(1)
	return StatusOk
}

func (thread *Thread) Get(slug string, id int) Status {
	var rows *sql.Rows
	var err error
	if id != 0 {
		rows, err = db.Query(sqlGetById, id)
	} else {
		rows, err = db.Query(sqlGetBySlug, slug)
	}
	defer rows.Close()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	if rows.Next() {
		slugNullable := &sql.NullString{}
		err = rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID, &thread.Votes)
		thread.Slug = slugNullable.String
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		return StatusOk
	}
	return StatusNotFound
}

func (update *ThreadUpdate) Update(slug string, id int) *Thread {
	thread := &Thread{}
	var err error
	slugNullable := &sql.NullString{}
	if id != 0 {
		err = db.QueryRow(sqlUpdateById, update.Message, update.Title, id).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID, &thread.Votes)
	} else {
		err = db.QueryRow(sqlUpdateBySlug, update.Message, update.Title, slug).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID, &thread.Votes)
	}
	thread.Slug = slugNullable.String
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	return thread
}

func GetByForumSlug(slug string, limit int, since *strfmt.DateTime, desc bool) (ThreadPointList, Status) {
	if forum.Get(slug) == nil {
		return nil, StatusUserOrForumNotExist
	}
	threads := make(ThreadPointList, 0, limit)
	var limitCondition string
	var sinceCondition string
	var orderCondition string
	var rows *sql.Rows
	var err error
	if desc {
		orderCondition = "ORDER BY created DESC"
	} else {
		orderCondition = "ORDER BY created"
	}
	if limit > 0 {
		limitCondition = "LIMIT $2"
	} else {
		limitCondition = ""
	}
	if since != nil {
		if desc == true {
			sinceCondition = "AND created <= $3"
		} else {
			sinceCondition = "AND created >= $3"
		}
		rows, err = db.Query(fmt.Sprintf("%s %s %s %s", sqlGetByForumSlug, sinceCondition, orderCondition, limitCondition), slug, limit, since)
	} else {
		rows, err = db.Query(fmt.Sprintf("%s %s %s", sqlGetByForumSlug, orderCondition, limitCondition), slug, limit)
	}

	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	slugNullable := &sql.NullString{}
	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID, &thread.Votes)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		thread.Slug = slugNullable.String
		threads = append(threads, thread)
	}
	rows.Close()
	return threads, StatusOk
}

func VoteForThread(slug string, id int, vote *Vote) *Thread {
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return nil
	}
	prevVoice := &sql.NullInt64{}
	threadId := &sql.NullInt64{}
	threadVotes := &sql.NullInt64{}
	userNickname := &sql.NullString{}
	if id != 0 {
		err = tx.QueryRow(sqlSelectThreadAndVoteById, id, vote.Nickname).Scan(prevVoice, threadId, &threadVotes, userNickname)
	} else {
		err = tx.QueryRow(sqlSelectThreadAndVoteBySlug, slug, vote.Nickname).Scan(prevVoice, threadId, &threadVotes, userNickname)
	}
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return nil
	}
	if !threadId.Valid || !userNickname.Valid {
		return nil
	}
	var prevVoiceInt int64
	if prevVoice.Valid {
		prevVoiceInt = prevVoice.Int64
		_, err = tx.Exec(sqlUpdateVote, threadId.Int64, userNickname.String, vote.Voice)
	} else {
		_, err = tx.Exec(sqlInsertVote, threadId.Int64, userNickname.String, vote.Voice)
	}
	newVotes := threadVotes.Int64 + (int64(vote.Voice) - prevVoiceInt)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return nil
	}
	thread := &Thread{}
	slugNullable := &sql.NullString{}
	err = tx.QueryRow(sqlUpdateThreadVotes, newVotes, threadId.Int64).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID, &thread.Votes)
	thread.Slug = slugNullable.String
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return nil
	}
	tx.Commit()
	return thread
}

func GetThreadId(slug string, id int) int {
	if id == 0 {
		err := db.QueryRow(sqlGetThreadIdBySlug, slug).Scan(&id)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
			return 0
		}
	} else {
		err := db.QueryRow(sqlGetThreadIdById, id).Scan(&id)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
			return 0
		}
	}
	return id
}
