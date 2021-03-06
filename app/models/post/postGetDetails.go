package post

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"tech-db-server/app/models/forum"
	"tech-db-server/app/models/thread"
	"tech-db-server/app/models/user"
)

const sqlGetPostData = `
	SELECT p.author, p.created, p.forum, p.isedited, p.message, p.parent, p.thread, %s, %s, %s
	FROM posts p
	LEFT JOIN threads t ON p.thread = t.id AND TRUE = $1
	LEFT JOIN forums f ON p.forum = f.slug AND TRUE = $2
	LEFT JOIN users u ON p.author = u.nickname AND TRUE = $3
	WHERE p.id = $4
`
const sqlThreadFields = `
t.author, t.created, t.forum, t."message", t.slug, t.title, t.id, t.votes
`
const sqlThreadFieldsEmpty = `
'', NULL, '', '', '', '', 0, 0
`
const sqlForumFields = `
f.slug, f."user", f.title, f.threads, f.posts
`
const sqlForumFieldsEmpty = `
'', '', '', 0, 0
`
const sqlUserFields = `
u.about, u.email, u.fullname, u.nickname
`
const sqlUserFieldsEmpty = `
'', '', '', ''
`

func PostDetails(id int32, related []string) *PostFull {
	post := &Post{}
	post.ID = id
	t := &thread.Thread{}
	f := &forum.Forum{}
	u := &user.User{}
	getAuthor, getForum, getThread := parseFlagsFromRelated(related)
	slugNullable := &pgtype.Varchar{}
	threadCreatedNullable := pgtype.Timestamptz{}
	err := db.QueryRow(*buildPostDetailsQuery(getAuthor, getForum, getThread), getThread, getForum, getAuthor, id).
		Scan(&post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread,
			&t.Author, &threadCreatedNullable, &t.Forum, &t.Message, slugNullable, &t.Title, &t.ID, &t.Votes,
			&f.Slug, &f.User, &f.Title, &f.Threads, &f.Posts,
			&u.About, &u.Email, &u.Fullname, &u.Nickname)
	t.Slug = slugNullable.String
	t.Created = threadCreatedNullable.Time
	if err == pgx.ErrNoRows {
		return nil
	}
	data := &PostFull{Post: post}
	if getAuthor {
		data.Author = u
	}
	if getForum {
		data.Forum = f
	}
	if getThread {
		data.Thread = t
	}
	return data
}

func parseFlagsFromRelated(related []string) (getAuthor bool, getForum bool, getThread bool) {
	for _, data := range related {
		switch data {
		case "forum":
			getForum = true
		case "thread":
			getThread = true
		case "user":
			getAuthor = true
		}
	}
	return
}

func buildPostDetailsQuery(getAuthor bool, getForum bool, getThread bool) *string {
	var forums string
	var threads string
	var users string
	if getAuthor {
		users = sqlUserFields
	} else {
		users = sqlUserFieldsEmpty
	}
	if getForum {
		forums = sqlForumFields
	} else {
		forums = sqlForumFieldsEmpty
	}
	if getThread {
		threads = sqlThreadFields
	} else {
		threads = sqlThreadFieldsEmpty
	}
	query := fmt.Sprintf(sqlGetPostData, threads, forums, users)
	return &query
}
