package repos

import (
	"database/sql"
	"friendfy-api/src/models"
)

type Posts struct {
	db *sql.DB
}

func NewRepositoryOfPosts(db *sql.DB) *Posts {
	return &Posts{db: db}
}

func (rep Posts) Create(post models.Post) (uint64, error) {
	statement, err := rep.db.Prepare(
		"INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}
