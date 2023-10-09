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

func (rep Posts) Get(userID uint64) ([]models.Post, error) {
	rows, err := rep.db.Query(`
	SELECT DISTINCT p.*, u.nick FROM posts p 
	INNER JOIN users u ON u.id = p.author_id 
	INNER JOIN followers f ON p.author_id = f.user_id
	WHERE u.id = ? OR f.follower_id = ?`,
		userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (rep Posts) GetById(id uint64) (models.Post, error) {
	row, err := rep.db.Query(`
	SELECT p.*, u.nick from
	posts p INNER JOIN users u
	on u.id = p.author_id where p.id = ? `,
		id,
	)
	if err != nil {
		return models.Post{}, err
	}
	defer row.Close()

	var post models.Post
	if row.Next() {
		err = row.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		)
		if err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}

func (rep Posts) Update(postID uint64, post models.Post) error {
	statement, err := rep.db.Prepare("UPDATE posts SET title = ?, content = ? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(post.Title, post.Content, postID)
	if err != nil {
		return err
	}

	return nil
}

func (rep Posts) Delete(postID uint64) error {
	statement, err := rep.db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(postID); err != nil {
		return err
	}

	return nil
}
