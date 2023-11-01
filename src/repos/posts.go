package repos

import (
	"database/sql"
	"errors"
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

func (rep Posts) GetThisUserPosts(userID uint64) ([]models.Post, error) {
	rows, err := rep.db.Query(
		`SELECT p.*, u.nick FROM posts p
		JOIN users u ON u.id = p.author_id
		WHERE p.author_id = ?`,
		userID,
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

func (rep Posts) Like(userID, postID uint64) error {
	row := rep.db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID)
	var authorID uint64
	err := row.Scan(&authorID)
	if err != nil {
		return err
	}
	if authorID == userID {
		return errors.New("user cannot like their own post")
	}

	row = rep.db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already liked this post")
	}

	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO post_likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rep Posts) Dislike(userID, postID uint64) error {
	row := rep.db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID)
	var authorID uint64
	err := row.Scan(&authorID)
	if err != nil {
		return err
	}
	if authorID == userID {
		return errors.New("user cannot dislike their own post")
	}

	row = rep.db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("there are no likes from you for this post")
	}

	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM post_likes WHERE user_id =? AND post_id = ?", userID, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
