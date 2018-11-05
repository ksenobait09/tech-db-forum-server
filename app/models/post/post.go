package post

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/lib/pq"
	"strconv"
	"strings"
	"tech-db-server/app/database"
	"tech-db-server/app/singletoneLogger"
	"tech-db-server/app/models/user"
	forumModel "tech-db-server/app/models/forum"
	"tech-db-server/app/models/thread"
)

var db *sql.DB

func init() {
	db = database.GetInstance()
}

//easyjson:json
type PostFull struct {
	Author *user.User `json:"author,omitempty"`
	Forum *forumModel.Forum `json:"forum,omitempty"`
	Post *Post `json:"post,omitempty"`
	Thread *thread.Thread `json:"thread,omitempty"`
}

//easyjson:json
type PostPointList []*Post

//easyjson:json
type Post struct {
	// Автор, написавший данное сообщение.
	// Required: true
	Author string `json:"author"`

	// Дата создания сообщения на форуме.
	// Read Only: true
	// Format: date-time
	Created *strfmt.DateTime `json:"created,omitempty"`

	// Идентификатор форума (slug) данного сообещния.
	// Read Only: true
	Forum string `json:"forum,omitempty"`

	// Идентификатор данного сообщения.
	// Read Only: true
	ID int64 `json:"id,omitempty"`

	// Истина, если данное сообщение было изменено.
	// Read Only: true
	IsEdited bool `json:"isEdited,omitempty"`

	// Собственно сообщение форума.
	// Required: true
	Message string `json:"message"`

	// Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	//
	Parent int64 `json:"parent,omitempty"`

	// Идентификатор ветви (id) обсуждения данного сообещния.
	// Read Only: true
	Thread int `json:"thread,omitempty"`

	Path pq.Int64Array `json:"-"`

	RootParent int64 `json:"-"`
}

type Status int

const (
	StatusConflict Status = iota + 1
	StatusNotExist
	StatusBadParent
	StatusOK
)

const sqlGetNextIds = `
	SELECT nextval(pg_get_serial_sequence('posts', 'id'))
	FROM generate_series(1, %d);
`

const sqlGetIdAndPathOfPostsInThread = `
	SELECT id, path FROM posts
	WHERE id IN (%s)
	AND thread = $1
`
const sqlUpdatePostsCount = `
	UPDATE forums
	SET posts = posts + $1
	WHERE slug = $2
`
const sqlGetThreadIdAndForumBySlug = `
	SELECT id, forum
	FROM threads
	WHERE slug = $1
`
const sqlGetThreadIdAndForumById = `
	SELECT id, forum
	FROM threads
	WHERE id = $1
`
const sqlGetCreatedFromPost = `
	SELECT created
	FROM posts
	WHERE id = $1
`

const sqlUpdate = `
	UPDATE posts SET isedited = message <> $1, message = $2 
	WHERE id = $3
	RETURNING author, created, forum, parent, thread, isedited
`

const sqlGet = `
	SELECT p.author, p.created, p.forum, p.isedited, p.message, p.parent, p.thread
	FROM posts p
	WHERE p.id = $1
`

const sqlGetPostsFlat = `
	SELECT p.author, p.created, p.forum, p.isedited, p.message, p.parent, p.thread, p.id
	FROM posts p
	WHERE thread = $1
`

const sqlGetPostsParentTree = `
	SELECT p.author, p.created, p.forum, p.isedited, p.message, p.parent, p.thread, p.id
	FROM posts p
	WHERE rootparent IN (SELECT id FROM posts p2 WHERE p2.thread=$1 AND p2.parent=0
`

func CreatePosts(threadSlug string, threadId int, posts PostPointList) (Status, PostPointList) {
	postsLen := len(posts)
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	var forum string
	//Проверка существования thread
	if threadId == 0 {
		err = tx.QueryRow(sqlGetThreadIdAndForumBySlug, threadSlug).Scan(&threadId, &forum)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
			return StatusNotExist, nil
		}
	} else {
		err = tx.QueryRow(sqlGetThreadIdAndForumById, threadId).Scan(&threadId, &forum)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
			return StatusNotExist, nil
		}
	}
	if postsLen == 0 {
		return StatusOK, posts
	}
	// Взять path родительских постов и authormap
	mapOfParentPathsById := make(map[int64][]int64)
	mapOfAuthors := make(map[string]bool)
	for _, post := range posts {
		if post.Parent != 0 {
			mapOfParentPathsById[post.Parent] = nil
		}
		mapOfAuthors[post.Author] = true
	}

	var parentIds []string
	for parentId := range mapOfParentPathsById {
		parentIds = append(parentIds, strconv.FormatInt(parentId, 10))
	}
	if len(parentIds) > 0 {
		returnedPostsCount := 0
		rows, err := tx.Query(fmt.Sprintf(sqlGetIdAndPathOfPostsInThread, strings.Join(parentIds, ", ")), threadId)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		for rows.Next() {
			returnedPostsCount++
			var id int64
			var path []int64
			err = rows.Scan(&id, pq.Array(&path))
			if err != nil {
				singletoneLogger.LogErrorWithStack(err)
			}
			mapOfParentPathsById[id] = path
		}
		rows.Close()

		// хотя бы один родительский пост не в той ветке
		if returnedPostsCount != len(mapOfParentPathsById) {
			return StatusBadParent, nil
		}
	}

	//взятие id для постов
	postIdsRows, err := tx.Query(fmt.Sprintf(sqlGetNextIds, postsLen))
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusNotExist, nil
	}
	var postIds []int64
	for postIdsRows.Next() {
		var availableId int64
		err = postIdsRows.Scan(&availableId)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		postIds = append(postIds, availableId)
	}
	postIdsRows.Close()

	// сохранение постов
	stmt, err := tx.Prepare(pq.CopyIn("posts", "id", "author", "forum", "message", "parent", "path", "rootparent", "thread"))
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	for i, post := range posts {
		post.ID = postIds[i]
		post.Forum = forum
		post.Thread = threadId
		if post.Parent > 0 {
			post.Path = append(mapOfParentPathsById[post.Parent], post.ID)
			post.RootParent = post.Path[0]
		} else {
			post.Path = append(post.Path, post.ID)
			post.RootParent = post.ID
		}
		_, err = stmt.Exec(post.ID, post.Author, post.Forum, post.Message, post.Parent, pq.Array(post.Path), post.RootParent, post.Thread)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
			return StatusNotExist, nil
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusNotExist, nil
	}
	err = stmt.Close()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusNotExist, nil
	}
	_, err = tx.Exec(sqlUpdatePostsCount, postsLen, forum)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusNotExist, nil
	}
	// userforum
	forumModel.InsertMapIntoUserForum(tx, forum, mapOfAuthors)

	err = tx.Commit()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusNotExist, nil
	}
	// Взять время создания постов
	created := strfmt.NewDateTime()
	err = db.QueryRow(sqlGetCreatedFromPost, posts[0].ID).Scan(&created)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
		return StatusNotExist, nil
	}
	for _, post := range posts {
		post.Created = &created
	}
	return StatusOK, posts
}

func (post *Post) Update() Status {
	if post.Message == "" {
		err := db.QueryRow(sqlGet, post.ID).
			Scan(&post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		if err == sql.ErrNoRows {
			return StatusNotExist
		}
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		return StatusOK
	}
	err := db.QueryRow(sqlUpdate, post.Message, post.Message, post.ID).
		Scan(&post.Author, &post.Created, &post.Forum, &post.Parent, &post.Thread, &post.IsEdited)
	if err == sql.ErrNoRows {
		return StatusNotExist
	}
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	return StatusOK
}

func GetPosts(id int, limit int, since int, sort string, desc bool) (PostPointList) {
	if limit == 0 {
		limit = 100
	}
	posts := make(PostPointList, 0, limit)
	var query strings.Builder
	query.Grow(100)

	switch sort {
	case "":
		fallthrough
	case "flat":
		fmt.Fprint(&query, sqlGetPostsFlat)
		if since != 0 {
			if desc {
				fmt.Fprint(&query, " AND p.id < $2")
			} else {
				fmt.Fprint(&query, " AND p.id > $2")
			}
		} else {
			fmt.Fprint(&query, " AND $2 = 0")
		}
		if desc {
			fmt.Fprint(&query, " ORDER BY p.id DESC")
		} else {
			fmt.Fprint(&query, " ORDER BY p.id")
		}
		fmt.Fprint(&query, " LIMIT $3")
	case "tree":
		fmt.Fprint(&query, sqlGetPostsFlat)
		if since != 0 {
			if desc {
				fmt.Fprint(&query, " AND p.path < (SELECT path FROM posts WHERE id = $2)")
			} else {
				fmt.Fprint(&query, " AND p.path > (SELECT path FROM posts WHERE id = $2)")
			}
		} else {
			fmt.Fprint(&query, " AND $2 = 0")
		}
		if desc {
			fmt.Fprint(&query, " ORDER BY p.path DESC")
		} else {
			fmt.Fprint(&query, " ORDER BY p.path")
		}
		fmt.Fprint(&query, " LIMIT $3")
	case "parent_tree":
		fmt.Fprint(&query, sqlGetPostsParentTree)
		if since != 0 {
			if desc {
				fmt.Fprint(&query, " AND p2.id < (SELECT rootparent FROM posts WHERE id=$2)")
			} else {
				fmt.Fprint(&query, " AND p2.id > (SELECT rootparent FROM posts WHERE id=$2)")
			}
		} else {
			fmt.Fprint(&query, " AND $2 = 0")
		}
		if desc {
			fmt.Fprint(&query, " ORDER BY p2.id DESC")
		} else {
			fmt.Fprint(&query, " ORDER BY p2.id")
		}
		fmt.Fprint(&query, " LIMIT $3)")
		if desc {
			fmt.Fprint(&query, " ORDER BY p.rootparent DESC, p.path")
		} else {
			fmt.Fprint(&query, " ORDER BY p.path")
		}
	}
	rows, err := db.Query(query.String(), id, since, limit)
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.ID)
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts

}