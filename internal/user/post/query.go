package post

import (
	"context"
	"time"
)

const (
	PostsByUserId int = iota
	InsertPost
	GetUserFeedById
)

type Query struct {
	SQL     string
	Timeout time.Duration
}

func GetQuery(queryIndex int) (string, context.Context, context.CancelFunc) {
	context, cancel := context.WithTimeout(context.Background(), queryMap[queryIndex].Timeout)
	return queryMap[queryIndex].SQL, context, cancel
}

var queryMap map[int]Query

func init() {
	queryMap = make(map[int]Query)

	queryMap[PostsByUserId] = Query{
		SQL: `SELECT p.id
					, p.created_at
					, p.updated_at
					, p.body
					, u.first_name
					, u.last_name
					, u.login
					, u.id
			  FROM posts as p
			  LEFT JOIN users u on u.id = p.user_id
			  WHERE p.user_id = ?
			  ORDER BY p.created_at desc`,
		Timeout: time.Second * 10,
	}

	queryMap[InsertPost] = Query{
		SQL: `INSERT INTO posts (user_id, body) 
			  VALUES (?, ?);`,
		Timeout: time.Second * 10,
	}

	queryMap[GetUserFeedById] = Query{
		SQL: `SELECT p.id
					, p.created_at
					, p.updated_at
					, p.body
					, u.first_name
					, u.last_name
					, u.login
					, u.id
              FROM posts p
			  LEFT JOIN users u on u.id = p.user_id
              WHERE p.user_id IN (
              	SELECT f.friend_id
              	FROM friends f
              	WHERE user_id = ?
              )
              ORDER BY p.created_at desc
			  LIMIT 1000`,
		Timeout: time.Second * 40,
	}
}
